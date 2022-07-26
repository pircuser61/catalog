package cache

import "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"

type Interface interface {
	Add(models.Good) error
	Get(uint64) (*models.Good, error)
	Update(models.Good) error
	Delete(uint64) error
	List() []models.Good
	//GetNextCode() uint64
}
