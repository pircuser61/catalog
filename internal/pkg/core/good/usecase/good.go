package usecase

import (
	"context"

	"github.com/pkg/errors"
	countryPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	unitOfMeasurePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type GoodUseCase struct {
	repository        goodPkg.Repository
	uomRepository     unitOfMeasurePkg.Repository
	countryRepository countryPkg.Repository
}

func New(repository goodPkg.Repository,
	uomRepository unitOfMeasurePkg.Repository,
	countryRepository countryPkg.Repository) goodPkg.Interface {
	return &GoodUseCase{
		repository:        repository,
		countryRepository: countryRepository,
		uomRepository:     uomRepository,
	}
}

func (c *GoodUseCase) Add(ctx context.Context, g *models.Good) error {

	uom_id, country_id, err := c.GetFieldsId(ctx, g)
	if err != nil {
		return err
	}

	return c.repository.Add(ctx, g.Name, uom_id, country_id)
}

func (c *GoodUseCase) Get(ctx context.Context, code uint64) (*models.Good, error) {
	result, err := c.repository.Get(ctx, code)
	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return nil, goodPkg.ErrGoodNotFound
	}
	return result, err
}

func (c *GoodUseCase) Update(ctx context.Context, g *models.Good) error {

	uom_id, country_id, err := c.GetFieldsId(ctx, g)
	if err != nil {
		return err
	}

	err = c.repository.Update(ctx, g.Code, g.Name, uom_id, country_id)
	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return goodPkg.ErrGoodNotFound
	}
	return err
}

func (c *GoodUseCase) Delete(ctx context.Context, code uint64) error {
	err := c.repository.Delete(ctx, code)

	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return goodPkg.ErrGoodNotFound
	}
	return err
}

func (c *GoodUseCase) List(ctx context.Context) ([]*models.Good, error) {
	return c.repository.List(ctx)
}

func (c *GoodUseCase) GetFieldsId(ctx context.Context, g *models.Good) (uint32, uint32, error) {
	if err := g.Validate(); err != nil {
		return 0, 0, err
	}
	uom, err := c.uomRepository.GetByName(ctx, g.UnitOfMeasure)
	if err != nil {
		if errors.Is(err, storePkg.ErrNotExists) {
			return 0, 0, errors.WithMessagef(models.ErrValidation,
				"Единица измерения %s не найдена в справочнике",
				g.UnitOfMeasure)
		}
		return 0, 0, err
	}
	country, err := c.countryRepository.GetByNane(ctx, g.Country)
	if err != nil {
		if errors.Is(err, storePkg.ErrNotExists) {
			return 0, 0, errors.WithMessagef(models.ErrValidation,
				"Страна %s не найдена в справочнике",
				g.Country)
		}
		return 0, 0, err
	}
	return uom.UnitOfMeasureId, country.CountryId, nil
}
