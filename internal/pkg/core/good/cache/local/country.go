package local

import (
	"context"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

func (c *cache) CountryAdd(context.Context, *models.Country) error {
	return errNe
}
func (c *cache) CountryGet(context.Context, uint32) (*models.Country, error) {
	return nil, errNe
}
func (c *cache) CountryUpdate(context.Context, *models.Country) error {
	return errNe
}
func (c *cache) CountryDelete(context.Context, uint32) error {
	return errNe
}
func (c *cache) CountryList(context.Context) ([]*models.Country, error) {
	resutl := make([]*models.Country, 1)
	return resutl, errNe
}
