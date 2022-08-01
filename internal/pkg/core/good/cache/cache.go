package cache

import (
	"context"
	"errors"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

var (
	ErrUserNotExists = errors.New("good does not exist")
	ErrUserExists    = errors.New("good exist")
	ErrTimeout       = errors.New("Timeout")
)

type Interface interface {
	Add(context.Context, *models.Good) error
	Get(context.Context, uint64) (*models.Good, error)
	Update(context.Context, *models.Good) error
	Delete(context.Context, uint64) error
	List(context.Context) ([]*models.Good, error)
	Disconnect(context.Context) error

	/* Методы для тестирования */
	Lock() string
	RLock() string
	Unlock() string
	RUnlock() string
}
