package usecase

import (
	"context"
	"errors"

	unitOfMeasurePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type UnitOfMeasureUseCase struct {
	repository unitOfMeasurePkg.Repository
}

func New(repository unitOfMeasurePkg.Repository) unitOfMeasurePkg.Interface {
	return &UnitOfMeasureUseCase{
		repository: repository,
	}
}

func (c *UnitOfMeasureUseCase) Add(ctx context.Context, uom *models.UnitOfMeasure) error {
	return c.repository.Add(ctx, uom)
}
func (c *UnitOfMeasureUseCase) Update(ctx context.Context, uom *models.UnitOfMeasure) error {
	err := c.repository.Update(ctx, uom)
	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return unitOfMeasurePkg.ErrUnitOfMeasurePkgNotFound
	}
	return err
}
func (c *UnitOfMeasureUseCase) Delete(ctx context.Context, unit_of_measure_id uint32) error {
	err := c.repository.Delete(ctx, unit_of_measure_id)
	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return unitOfMeasurePkg.ErrUnitOfMeasurePkgNotFound
	}
	return err
}
func (c *UnitOfMeasureUseCase) Get(ctx context.Context, unit_of_measure_id uint32) (*models.UnitOfMeasure, error) {
	result, err := c.repository.Get(ctx, unit_of_measure_id)
	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return nil, unitOfMeasurePkg.ErrUnitOfMeasurePkgNotFound
	}
	return result, err
}

func (c *UnitOfMeasureUseCase) List(ctx context.Context) ([]*models.UnitOfMeasure, error) {
	return c.repository.List(ctx)
}
