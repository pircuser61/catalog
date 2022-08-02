package local

import (
	"context"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

func (c *cache) UnitOfMeasureAdd(context.Context, *models.UnitOfMeasure) error {
	return errNe
}
func (c *cache) UnitOfMeasureGet(context.Context, uint32) (*models.UnitOfMeasure, error) {
	return nil, errNe
}
func (c *cache) UnitOfMeasureUpdate(context.Context, *models.UnitOfMeasure) error {
	return errNe
}
func (c *cache) UnitOfMeasureDelete(context.Context, uint32) error {
	return errNe
}
func (c *cache) UnitOfMeasureList(context.Context) ([]*models.UnitOfMeasure, error) {
	result := make([]*models.UnitOfMeasure, 1)
	return result, errNe
}
