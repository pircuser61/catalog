package good

import (
	"errors"

	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	localCachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache/local"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

type Interface interface {
	Create(models.Good) error
	Update(models.Good) error
	Delete(uint64) error
	Get(uint64) (*models.Good, error)
	List() []models.Good
}

type core struct {
	cache cachePkg.Interface
}

var ErrValidation = errors.New("Invalid data")

func New() Interface {
	return &core{
		cache: localCachePkg.New(),
	}
}

func (c *core) Create(g models.Good) error {
	if err := g.Validate(); err != nil {
		return err
	}
	return c.cache.Add(g)
}

func (c *core) Get(code uint64) (*models.Good, error) {
	return c.cache.Get(code)
}

func (c *core) Update(g models.Good) error {
	if err := g.Validate(); err != nil {
		return err
	}
	return c.cache.Update(g)
}

func (c *core) Delete(code uint64) error {
	return c.cache.Delete(code)
}

func (c *core) List() []models.Good {
	return c.cache.List()
}
