package cache

import (
	"errors"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

var (
	ErrUserNotExists = errors.New("good does not exist")
	ErrUserExists    = errors.New("good exist")
	ErrTimeout       = errors.New("Timeout")
)

type Interface interface {
	Add(models.Good) error
	Get(uint64) (*models.Good, error)
	Update(models.Good) error
	Delete(uint64) error
	List() ([]models.Good, error)

	Lock() string
	RLock() string
	Unlock() string
	RUnlock() string
}
