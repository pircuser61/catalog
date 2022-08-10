package models

import (
	"github.com/pkg/errors"
)

var ErrValidation = errors.New("invalid data")

type UnitOfMeasure struct {
	Name            string
	UnitOfMeasureId uint32
}

type Country struct {
	Name      string
	CountryId uint32
}

type Good struct {
	Code          uint64
	Name          string
	UnitOfMeasure string
	Country       string
}
