package usecase

import (
	"context"
	"errors"

	countryPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type CountryUseCase struct {
	repository countryPkg.Repository
}

func New(repository countryPkg.Repository) countryPkg.Interface {
	return &CountryUseCase{
		repository: repository,
	}
}

func (c *CountryUseCase) Add(ctx context.Context, ct *models.Country) error {
	if err := ct.Validate(); err != nil {
		return err
	}
	return c.repository.Add(ctx, ct)
}

func (c *CountryUseCase) Update(ctx context.Context, ct *models.Country) error {
	err := ct.Validate()
	if err != nil {
		return err
	}
	err = c.repository.Update(ctx, ct)
	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return countryPkg.ErrCountryNotFound
	}
	return err
}

func (c *CountryUseCase) Delete(ctx context.Context, country_id uint32) error {
	err := c.repository.Delete(ctx, country_id)
	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return countryPkg.ErrCountryNotFound
	}
	return err
}

func (c *CountryUseCase) Get(ctx context.Context, country_id uint32) (*models.Country, error) {
	result, err := c.repository.Get(ctx, country_id)
	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return nil, countryPkg.ErrCountryNotFound
	}
	return result, err
}

func (c *CountryUseCase) List(ctx context.Context) ([]*models.Country, error) {
	return c.repository.List(ctx)
}

func (c *CountryUseCase) GetByName(ctx context.Context, countryName string) (*models.Country, error) {
	result, err := c.repository.GetByName(ctx, countryName)
	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return nil, countryPkg.ErrCountryNotFound
	}
	return result, err
}
