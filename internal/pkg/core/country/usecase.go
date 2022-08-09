package country

import (
	"context"
	"errors"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
)

type Interface interface {
	Add(context.Context, *models.Country) error
	Get(context.Context, uint32) (*models.Country, error)
	Update(context.Context, *models.Country) error
	Delete(context.Context, uint32) error
	List(context.Context) ([]*models.Country, error)
	GetByName(context.Context, string) (*models.Country, error)
}

var ErrCountryNotFound = errors.New("country not found")
