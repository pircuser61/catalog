package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/go-redis/redis"
	"gitlab.ozon.dev/pircuser61/catalog/config"
	counters "gitlab.ozon.dev/pircuser61/catalog/internal/counters"
	logger "gitlab.ozon.dev/pircuser61/catalog/internal/logger"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	goodRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/repository/postgre"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	"go.uber.org/zap"
)

type GoodCachedRepository struct {
	dbRepository goodPkg.Repository
	goodsRedis   *redis.Client
}

func New(pool pgxpoolmock.PgxPool, timeout time.Duration) goodPkg.Repository {

	goods := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		DB:       config.RedisGoodDb,
		Password: config.RedisPassword})

	return &GoodCachedRepository{
		dbRepository: goodRepo.New(pool, timeout),
		goodsRedis:   goods,
	}
}

func (c *GoodCachedRepository) setCache(good *models.Good) {
	data, err := json.Marshal(good)
	if err != nil {
		logger.Error("Cache convert to json", zap.Error(err))
		return
	}
	status := c.goodsRedis.Set(fmt.Sprint(good.Code), data, config.RedisExpiration)
	if status.Err() != nil {
		logger.Error("Cache set", zap.Error(err))
	}
}

func (c *GoodCachedRepository) Add(ctx context.Context, good *models.Good) error {
	return c.dbRepository.Add(ctx, good)
}

func (c *GoodCachedRepository) Update(ctx context.Context, good *models.Good) error {
	err := c.dbRepository.Update(ctx, good)
	if err == nil {
		c.setCache(good)
	}
	return err
}

func (c *GoodCachedRepository) Get(ctx context.Context, code uint64) (*models.Good, error) {
	var cacheNotFound = false
	key := fmt.Sprint(code)
	result := c.goodsRedis.Get(key)
	if err := result.Err(); err == redis.Nil {
		cacheNotFound = true
	} else if err != nil {
		logger.Error("Cache get", zap.Error(err))
	} else {
		counters.Hit()
		bytes, err := result.Bytes()
		if err == nil {
			good := &models.Good{}
			err = json.Unmarshal(bytes, good)
			if err == nil {
				return good, nil
			}
		}
		logger.Error("Cache parse", zap.Error(err))
	}

	good, err := c.dbRepository.Get(ctx, code)
	if err != nil {
		return nil, err
	}
	if cacheNotFound {
		counters.Miss()
	}
	c.setCache(good)
	return good, nil
}

func (c *GoodCachedRepository) Delete(ctx context.Context, code uint64) error {
	err := c.dbRepository.Delete(ctx, code)
	c.goodsRedis.Del(fmt.Sprint(code))
	return err
}

func (c *GoodCachedRepository) List(ctx context.Context, limit uint64, offset uint64) ([]*models.Good, error) {
	return c.dbRepository.List(ctx, limit, offset)
}

func (c *GoodCachedRepository) GetKeys(ctx context.Context, good *models.Good) (*goodPkg.GoodKeys, error) {
	return c.dbRepository.GetKeys(ctx, good)
}
