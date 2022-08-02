package good

import (
	"context"
	"errors"

	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

func (c *core) UnitOfMeasureCreate(ctx context.Context, uom models.UnitOfMeasure) error {
	if err := uom.Validate(); err != nil {
		return err
	}
	return c.cache.UnitOfMeasureAdd(ctx, &uom)
}
func (c *core) UnitOfMeasureUpdate(ctx context.Context, uom models.UnitOfMeasure) error {
	err := uom.Validate()
	if err != nil {
		return err
	}
	err = c.cache.UnitOfMeasureUpdate(ctx, &uom)
	if err != nil && errors.Is(err, cachePkg.ErrObjNotExists) {
		return ErrNotFound
	}
	return err
}
func (c *core) UnitOfMeasureDelete(ctx context.Context, unit_of_measure_id uint32) error {
	err := c.cache.UnitOfMeasureDelete(ctx, unit_of_measure_id)
	if err != nil && errors.Is(err, cachePkg.ErrObjNotExists) {
		return ErrNotFound
	}
	return err
}
func (c *core) UnitOfMeasureGet(ctx context.Context, unit_of_measure_id uint32) (*models.UnitOfMeasure, error) {
	result, err := c.cache.UnitOfMeasureGet(ctx, unit_of_measure_id)
	if err != nil && errors.Is(err, cachePkg.ErrObjNotExists) {
		return nil, ErrNotFound
	}
	return result, err
}
func (c *core) UnitOfMeasureList(ctx context.Context) ([]*models.UnitOfMeasure, error) {
	return c.cache.UnitOfMeasureList(ctx)
}
