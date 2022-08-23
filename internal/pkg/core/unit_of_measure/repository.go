//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository

package unit_of_measure

import (
	"context"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
)

type Repository interface {
	Add(context.Context, *models.UnitOfMeasure) error
	Get(context.Context, uint32) (*models.UnitOfMeasure, error)
	Update(context.Context, *models.UnitOfMeasure) error
	Delete(context.Context, uint32) error
	List(context.Context) ([]*models.UnitOfMeasure, error)
	GetByName(context.Context, string) (*models.UnitOfMeasure, error)
}
