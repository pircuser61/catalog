package unit_of_measure

import (
	"context"
	"errors"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
)

type Interface interface {
	Add(context.Context, *models.UnitOfMeasure) error
	Get(context.Context, uint32) (*models.UnitOfMeasure, error)
	Update(context.Context, *models.UnitOfMeasure) error
	Delete(context.Context, uint32) error
	List(context.Context) ([]*models.UnitOfMeasure, error)
}

var ErrUnitOfMeasurePkgNotFound = errors.New("UnitOfMeasurePkg not found")
