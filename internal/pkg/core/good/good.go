package good

import (
	"context"
	"errors"

	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

func (c *core) GoodCreate(ctx context.Context, g models.Good) error {
	if err := g.Validate(); err != nil {
		return err
	}
	return c.cache.GoodAdd(ctx, &g)
}

func (c *core) GoodGet(ctx context.Context, code uint64) (*models.Good, error) {
	result, err := c.cache.GoodGet(ctx, code)
	if err != nil && errors.Is(err, cachePkg.ErrObjNotExists) {
		return nil, ErrNotFound
	}
	return result, err
}

func (c *core) GoodUpdate(ctx context.Context, g models.Good) error {
	err := g.Validate()
	if err != nil {
		return err
	}
	err = c.cache.GoodUpdate(ctx, &g)
	if err != nil && errors.Is(err, cachePkg.ErrObjNotExists) {
		return ErrNotFound
	}
	return err
}

func (c *core) GoodDelete(ctx context.Context, code uint64) error {
	err := c.cache.GoodDelete(ctx, code)

	if err != nil && errors.Is(err, cachePkg.ErrObjNotExists) {
		return ErrNotFound
	}
	return err
}

func (c *core) GoodList(ctx context.Context) ([]*models.Good, error) {
	return c.cache.GoodList(ctx)
}
