package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-redis/redis"
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	"gitlab.ozon.dev/pircuser61/catalog/config"
	logger "gitlab.ozon.dev/pircuser61/catalog/internal/logger"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var redisClient *redis.Client

func init() {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		DB:       config.RedisResponseDb,
		Password: config.RedisPassword})
}

func isAsync(ctx context.Context) bool {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		strings := md.Get("mode")
		if len(strings) > 0 && strings[0] == "async" {
			return true
		}
	}
	return false
}

func (i *Implementation) GoodCreate(ctx context.Context, in *pb.GoodCreateRequest) (*pb.GoodCreateResponse, error) {

	good := &models.Good{
		Name:          in.GetName(),
		UnitOfMeasure: in.GetUnitOfMeasure(),
		Country:       in.GetCountry(),
	}
	if isAsync(ctx) {
		token := fmt.Sprintf("%v", time.Now())
		logger.Debug("Create async")
		go func() {
			ctx := context.Background()
			time.Sleep(time.Second * 2)
			err := i.good.Add(ctx, good)
			var text string
			if err == nil {
				text = "Good created"
			} else {
				text = "Error: " + err.Error()
			}

			redisClient.Set(token, text, config.RedisResponseExpiration)
		}()
		md := metadata.Pairs("token", token)
		grpc.SendHeader(ctx, md)
	} else {
		logger.Debug("Create ...")
		if err := i.good.Add(ctx, good); err != nil {
			if errors.Is(err, models.ErrValidation) {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			if errors.Is(err, storePkg.ErrTimeout) {
				return nil, status.Error(codes.DeadlineExceeded, err.Error())
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &pb.GoodCreateResponse{}, nil
}

func (i *Implementation) GoodGet(ctx context.Context, in *pb.GoodGetRequest) (*pb.GoodGetResponse, error) {
	good, err := i.good.Get(ctx, in.GetCode())
	if err != nil {
		if errors.Is(err, storePkg.ErrNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GoodGetResponse{
		Good: &pb.Good{
			Code:          good.Code,
			Name:          good.Name,
			UnitOfMeasure: good.UnitOfMeasure,
			Country:       good.Country},
	}, nil
}

func (i *Implementation) GoodUpdate(ctx context.Context, in *pb.GoodUpdateRequest) (*pb.GoodUpdateResponse, error) {
	inGood := in.GetGood()
	if inGood == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	good := &models.Good{
		Code:          inGood.GetCode(),
		Name:          inGood.GetName(),
		UnitOfMeasure: inGood.GetUnitOfMeasure(),
		Country:       inGood.GetCountry()}

	if isAsync(ctx) {
		go func() {
			ctx := context.Background()
			time.Sleep(time.Second * 2)
			err := i.good.Update(ctx, good)
			if err == nil {
				redisClient.Publish("response", "done")
			} else {
				redisClient.Publish("response", "Error: "+err.Error())
			}
		}()
	} else {

		if err := i.good.Update(ctx, good); err != nil {
			if errors.Is(err, models.ErrValidation) {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
			if errors.Is(err, storePkg.ErrNotExists) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			if errors.Is(err, storePkg.ErrTimeout) {
				return nil, status.Error(codes.DeadlineExceeded, err.Error())
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &pb.GoodUpdateResponse{}, nil
}

func (i *Implementation) GoodDelete(ctx context.Context, in *pb.GoodDeleteRequest) (*pb.GoodDeleteResponse, error) {
	if err := i.good.Delete(ctx, in.GetCode()); err != nil {
		if errors.Is(err, storePkg.ErrNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GoodDeleteResponse{}, nil
}

func (i *Implementation) GoodList(stream pb.Catalog_GoodListServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		result := getGoodList(i.good, in.GetLimit(), in.GetOffset())
		if err := stream.Send(result); err != nil {
			return err
		}
	}
}

func getGoodList(goodRepo good.Repository, limit, offset uint64) *pb.GoodListResponse {
	result := &pb.GoodListResponse{}
	ctx := context.Background()
	goods, err := goodRepo.List(ctx, limit, offset)
	if err != nil {
		if errors.Is(err, storePkg.ErrTimeout) {
			err = status.Error(codes.DeadlineExceeded, err.Error())
		} else {
			err = status.Error(codes.Internal, err.Error())
		}
		errMsg := err.Error()
		result.Error = &errMsg
	} else {
		listGoods := make([]*pb.GoodListResponse_Good, 0, len(goods))
		for _, good := range goods {
			listGoods = append(listGoods, &pb.GoodListResponse_Good{
				Code: good.Code,
				Name: good.Name,
			})
		}
		result.Goods = listGoods
	}
	return result
}
