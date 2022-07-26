package local

import (
	"github.com/pkg/errors"

	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

type cache struct {
	data     map[uint64]models.Good
	lastCode uint64
}

var (
	ErrUserNotExists = errors.New("good does not exist")
	ErrUserExists    = errors.New("good exist")
)

func New() cachePkg.Interface {
	return &cache{data: map[uint64]models.Good{}}
}

func (c *cache) GetNextCode() uint64 {
	c.lastCode++
	return c.lastCode
}

func (c *cache) List() []models.Good {
	result := make([]models.Good, 0, len(c.data))
	for _, x := range c.data {
		result = append(result, x)
	}
	return result
}

func (c *cache) Add(g models.Good) error {
	if _, ok := c.data[g.GetCode()]; ok {
		return errors.Wrapf(ErrUserExists, "code %d", g.GetCode())
	}
	if err := g.SetCode(c.GetNextCode()); err != nil {
		return err
	}
	c.data[g.GetCode()] = g
	return nil
}

func (c *cache) Get(code uint64) (*models.Good, error) {
	if g, ok := c.data[code]; ok {
		return &g, nil
	}
	return nil, errors.Wrapf(ErrUserNotExists, "code %d", code)
}

func (c *cache) Update(g models.Good) error {
	if _, ok := c.data[g.GetCode()]; !ok {
		return errors.Wrapf(ErrUserNotExists, "code %d", g.GetCode())
	}
	c.data[g.GetCode()] = g
	return nil
}

func (c *cache) Delete(code uint64) error {
	if _, ok := c.data[code]; ok {
		delete(c.data, code)
		return nil
	}
	return errors.Wrapf(ErrUserNotExists, "code %d", code)
}
