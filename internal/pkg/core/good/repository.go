package good

import (
	"context"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
)

type GoodKeys struct {
	TheOne          *int // место для "SELECT 1"
	UnitOfMeasureId *uint32
	CountryId       *uint32
}

type Repository interface {
	Add(context.Context, *models.Good, *GoodKeys) error
	Get(context.Context, uint64) (*models.Good, error)
	Update(context.Context, *models.Good, *GoodKeys) error
	Delete(context.Context, uint64) error
	GetKeys(context.Context, *models.Good) (*GoodKeys, error)
	List(context.Context, uint64, uint64) ([]*models.Good, error)
}
