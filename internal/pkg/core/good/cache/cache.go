package cache

import (
	"context"
	"errors"

	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

var (
	ErrObjNotExists = errors.New("obj does not exist")
	ErrObjExists    = errors.New("obj exist")
	ErrTimeout      = errors.New("Timeout")
)

type Interface interface {
	GoodAdd(context.Context, *models.Good) error
	GoodGet(context.Context, uint64) (*models.Good, error)
	GoodUpdate(context.Context, *models.Good) error
	GoodDelete(context.Context, uint64) error
	GoodList(context.Context) ([]*models.Good, error)

	CountryAdd(context.Context, *models.Country) error
	CountryGet(context.Context, uint32) (*models.Country, error)
	CountryUpdate(context.Context, *models.Country) error
	CountryDelete(context.Context, uint32) error
	CountryList(context.Context) ([]*models.Country, error)

	UnitOfMeasureAdd(context.Context, *models.UnitOfMeasure) error
	UnitOfMeasureGet(context.Context, uint32) (*models.UnitOfMeasure, error)
	UnitOfMeasureUpdate(context.Context, *models.UnitOfMeasure) error
	UnitOfMeasureDelete(context.Context, uint32) error
	UnitOfMeasureList(context.Context) ([]*models.UnitOfMeasure, error)

	Close(context.Context) error
}
