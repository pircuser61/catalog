package good

import (
	"context"
	"errors"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
)

type Interface interface {
	Add(context.Context, *models.Good) error
	Get(context.Context, uint64) (*models.Good, error)
	Update(context.Context, *models.Good) error
	Delete(context.Context, uint64) error
	List(context.Context, uint64, uint64) ([]*models.Good, error)
}

var ErrGoodNotFound = errors.New("good not found")
