package usecase

import (
	"context"

	"github.com/pkg/errors"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type GoodUseCase struct {
	repository goodPkg.Repository
}

func New(repository goodPkg.Repository) goodPkg.Interface {
	return &GoodUseCase{
		repository: repository,
	}
}

func (c *GoodUseCase) Add(ctx context.Context, good *models.Good) error {
	keys, err := c.repository.GetKeys(ctx, good)
	if err != nil {
		return err
	}

	return c.repository.Add(ctx, good, keys)
}

func (c *GoodUseCase) Update(ctx context.Context, good *models.Good) error {
	keys, err := c.repository.GetKeys(ctx, good)
	if err != nil {
		return err
	}

	return c.repository.Update(ctx, good, keys)
}

func (c *GoodUseCase) Get(ctx context.Context, code uint64) (*models.Good, error) {
	result, err := c.repository.Get(ctx, code)
	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return nil, goodPkg.ErrGoodNotFound
	}
	return result, err
}

func (c *GoodUseCase) Delete(ctx context.Context, code uint64) error {
	err := c.repository.Delete(ctx, code)

	if err != nil && errors.Is(err, storePkg.ErrNotExists) {
		return goodPkg.ErrGoodNotFound
	}
	return err
}

func (c *GoodUseCase) List(ctx context.Context, limit uint64, offset uint64) ([]*models.Good, error) {
	return c.repository.List(ctx, limit, offset)
}
