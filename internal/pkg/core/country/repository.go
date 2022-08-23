//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository

package country

import (
	"context"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
)

type Repository interface {
	Add(context.Context, *models.Country) error
	Get(context.Context, uint32) (*models.Country, error)
	Update(context.Context, *models.Country) error
	Delete(context.Context, uint32) error
	List(context.Context) ([]*models.Country, error)
	GetByName(context.Context, string) (*models.Country, error)
}
