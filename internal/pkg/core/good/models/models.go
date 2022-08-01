package models

import (
	"fmt"

	"github.com/pkg/errors"
)

var ErrValidation = errors.New("invalid data")

type Good struct {
	Code          uint64
	Name          string
	UnitOfMeasure string `db:"unit_of_measure"`
	Country       string
}

func (g *Good) GetCode() uint64 {
	return g.Code
}

func (g *Good) SetCode(code uint64) error {
	g.Code = code
	return nil
}

func (g Good) Validate() error {
	if len(g.Name) < 3 || len(g.Name) > 40 {
		return errors.WithMessagef(ErrValidation, "bad name <%v>", g.Name)
	}

	if len(g.UnitOfMeasure) > 10 {
		return errors.WithMessagef(ErrValidation, "bad unit of measure <%v>", g.UnitOfMeasure)
	}

	if len(g.Country) < 3 || len(g.Country) > 20 {
		return errors.WithMessagef(ErrValidation, "bad country <%v>", g.Country)
	}
	return nil
}

func (g *Good) String() string {
	return fmt.Sprintf("%d %s", g.GetCode(), g.GetName())
}

func (g *Good) GetName() string {
	return g.Name
}

func (g *Good) GetUnitOfMeasure() string {
	return g.UnitOfMeasure
}

func (g *Good) GetCountry() string {
	return g.Country
}
