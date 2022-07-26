package models

import "fmt"

type Good struct {
	Code          uint64
	Name          string
	UnitOfMeasure string
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
		return fmt.Errorf("bad name <%v>", g.Name)
	}

	if len(g.UnitOfMeasure) > 10 {
		return fmt.Errorf("bad unit of measure <%v>", g.UnitOfMeasure)
	}

	if len(g.Country) < 3 || len(g.Country) > 20 {
		return fmt.Errorf("bad country <%v>", g.Country)
	}
	return nil
}

func (g *Good) String() string {
	return fmt.Sprintf("%d %s (%s) %s", g.GetCode(), g.GetName(), g.GetUnitOfMeasure(), g.GetCountry())
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
