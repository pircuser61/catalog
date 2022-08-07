package repository

import (
	"errors"

	countryPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	unitOfMeasurePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure"
)

var (
	ErrNotExists = errors.New("obj does not exist")
	ErrExists    = errors.New("obj exist")
	ErrTimeout   = errors.New("Timeout")
)

type Core struct {
	Good          goodPkg.Interface
	Country       countryPkg.Interface
	UnitOfMeasure unitOfMeasurePkg.Interface
}
