package good

import (
	"context"
	"errors"

	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

func (c *core) CountryCreate(ctx context.Context, ct models.Country) error {
	if err := ct.Validate(); err != nil {
		return err
	}
	return c.cache.CountryAdd(ctx, &ct)
}
func (c *core) CountryUpdate(ctx context.Context, ct models.Country) error {
	err := ct.Validate()
	if err != nil {
		return err
	}
	err = c.cache.CountryUpdate(ctx, &ct)
	if err != nil && errors.Is(err, cachePkg.ErrObjNotExists) {
		return ErrNotFound
	}
	return err
}
func (c *core) CountryDelete(ctx context.Context, country_id uint32) error {
	err := c.cache.CountryDelete(ctx, country_id)
	if err != nil && errors.Is(err, cachePkg.ErrObjNotExists) {
		return ErrNotFound
	}
	return err
}
func (c *core) CountryGet(ctx context.Context, country_id uint32) (*models.Country, error) {
	result, err := c.cache.CountryGet(ctx, country_id)
	if err != nil && errors.Is(err, cachePkg.ErrObjNotExists) {
		return nil, ErrNotFound
	}
	return result, err
}
func (c *core) CountryList(ctx context.Context) ([]*models.Country, error) {
	return c.cache.CountryList(ctx)
}
