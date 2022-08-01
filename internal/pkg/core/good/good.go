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
	Create(context.Context, models.Good) error
	Update(context.Context, models.Good) error
	Delete(context.Context, uint64) error
	Get(context.Context, uint64) (*models.Good, error)
	List(context.Context) ([]*models.Good, error)
	GetCache() cachePkg.Interface
	Disconnect(context.Context) error
}

type core struct {
	cache cachePkg.Interface
}

var ErrNotFound = errors.New("good not found")

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

func (c *core) Create(ctx context.Context, g models.Good) error {
	if err := g.Validate(); err != nil {
		return err
	}
	return c.cache.Add(ctx, &g)
}

func (c *core) Get(ctx context.Context, code uint64) (*models.Good, error) {
	return c.cache.Get(ctx, code)
}

func (c *core) Update(ctx context.Context, g models.Good) error {
	err := g.Validate()
	if err != nil {
		return err
	}
	err = c.cache.Update(ctx, &g)
	if err != nil && errors.Is(err, cachePkg.ErrUserNotExists) {
		return ErrNotFound
	}
	return err
}

func (c *core) Delete(ctx context.Context, code uint64) error {
	err := c.cache.Delete(ctx, code)

	if err != nil && errors.Is(err, cachePkg.ErrUserNotExists) {
		return ErrNotFound
	}
	return err
}

func (c *core) List(ctx context.Context) ([]*models.Good, error) {
	return c.cache.List(ctx)
}

func (c *core) GetCache() cachePkg.Interface {
	return c.cache
}

func (c *core) Disconnect(ctx context.Context) error {
	return c.Disconnect(ctx)
}
