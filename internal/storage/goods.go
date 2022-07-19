package storage

import (
	"fmt"
)

var lastCode uint64

type Goods struct {
	code    uint64
	name    string
	uom     string
	country string
}

func NewGoods(name, uom, country string) (*Goods, error) {
	g := Goods{}
	if err := g.SetFields(name, uom, country); err != nil {
		return nil, err
	}
	lastCode++
	g.code = lastCode
	return &g, nil
}

func (g *Goods) SetFields(name, uom, country string) error {
	if err := g.SetName(name); err != nil {
		return err
	}

	if err := g.SetUom(uom); err != nil {
		return err
	}

	if err := g.SetCountry(country); err != nil {
		return err
	}
	return nil
}

func (g *Goods) String() string {
	return fmt.Sprintf("%d %s (%s) %s", g.GetCode(), g.GetName(), g.GetUom(), g.GetCountry())
}

func (g *Goods) SetName(name string) error {
	if len(name) < 3 || len(name) > 40 {
		return fmt.Errorf("bad name <%v>", name)
	}
	g.name = name
	return nil
}

func (g *Goods) SetUom(uom string) error {
	if len(uom) > 10 {
		return fmt.Errorf("bad unit of measure <%v>", uom)
	}
	g.uom = uom
	return nil
}

func (g *Goods) SetCountry(country string) error {
	if len(country) < 3 || len(country) > 20 {
		return fmt.Errorf("bad country <%v>", country)
	}
	g.country = country
	return nil
}

func (g *Goods) GetCode() uint64 {
	return g.code
}

func (g *Goods) GetName() string {
	return g.name
}

func (g *Goods) GetUom() string {
	return g.uom
}

func (g *Goods) GetCountry() string {
	return g.country
}
