package good

import (
	"context"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
)

type Repository interface {
	Add(context.Context, string, uint32, uint32) error
	Get(context.Context, uint64) (*models.Good, error)
	Update(context.Context, uint64, string, uint32, uint32) error
	Delete(context.Context, uint64) error
	List(context.Context) ([]*models.Good, error)
}
