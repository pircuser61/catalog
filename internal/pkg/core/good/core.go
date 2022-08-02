package good

import (
	"context"
	"errors"

	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	localCachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache/local"
	postgrePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache/postrgre"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

type Interface interface {
	GetCache() cachePkg.Interface
	Close(context.Context) error

	GoodCreate(context.Context, models.Good) error
	GoodUpdate(context.Context, models.Good) error
	GoodDelete(context.Context, uint64) error
	GoodGet(context.Context, uint64) (*models.Good, error)
	GoodList(context.Context) ([]*models.Good, error)

	CountryCreate(context.Context, models.Country) error
	CountryUpdate(context.Context, models.Country) error
	CountryDelete(context.Context, uint32) error
	CountryGet(context.Context, uint32) (*models.Country, error)
	CountryList(context.Context) ([]*models.Country, error)

	UnitOfMeasureCreate(context.Context, models.UnitOfMeasure) error
	UnitOfMeasureUpdate(context.Context, models.UnitOfMeasure) error
	UnitOfMeasureDelete(context.Context, uint32) error
	UnitOfMeasureGet(context.Context, uint32) (*models.UnitOfMeasure, error)
	UnitOfMeasureList(context.Context) ([]*models.UnitOfMeasure, error)
}

type core struct {
	cache cachePkg.Interface
}

var ErrNotFound = errors.New("not found")

func New() Interface {
	return &core{
		cache: localCachePkg.New(),
	}
}

func NewPostgre(ctx context.Context) Interface {
	return &core{
		cache: postgrePkg.New(ctx),
	}
}

func (c *core) GetCache() cachePkg.Interface {
	return c.cache
}

func (c *core) Close(ctx context.Context) error {
	return c.Close(ctx)
}
