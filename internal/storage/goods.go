package storage

import (
	"fmt"
)

var lastCode uint64

type Good struct {
	code          uint64
	name          string
	unitOfMeasure string
	country       string
}

func NewGood(name, uom, country string) (*Good, error) {
	g := Good{}
	if err := g.SetFields(name, uom, country); err != nil {
		return nil, err
	}
	lastCode++
	g.code = lastCode
	return &g, nil
}

func (g *Good) SetFields(name, uom, country string) error {
	if err := g.SetName(name); err != nil {
		return err
	}

	if err := g.SetUnitOfMeasure(uom); err != nil {
		return err
	}

	if err := g.SetCountry(country); err != nil {
		return err
	}
	return nil
}

func (g *Good) String() string {
	return fmt.Sprintf("%d %s (%s) %s", g.GetCode(), g.GetName(), g.GetUnitOfMeasure(), g.GetCountry())
}

func (g *Good) SetName(name string) error {
	if len(name) < 3 || len(name) > 40 {
		return fmt.Errorf("bad name <%v>", name)
	}
	g.name = name
	return nil
}

func (g *Good) SetUnitOfMeasure(uom string) error {
	if len(uom) > 10 {
		return fmt.Errorf("bad unit of measure <%v>", uom)
	}
	g.unitOfMeasure = uom
	return nil
}

func (g *Good) SetCountry(country string) error {
	if len(country) < 3 || len(country) > 20 {
		return fmt.Errorf("bad country <%v>", country)
	}
	g.country = country
	return nil
}

func (g *Good) GetCode() uint64 {
	return g.code
}

func (g *Good) GetName() string {
	return g.name
}

func (g *Good) GetUnitOfMeasure() string {
	return g.unitOfMeasure
}

func (g *Good) GetCountry() string {
	return g.country
}
