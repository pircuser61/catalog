package fixtures

import models "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"

type GoodBuilder struct {
	instance *models.Good
}

func Good() *GoodBuilder {
	return &GoodBuilder{instance: &models.Good{}}
}

func (gb GoodBuilder) Code(code uint64) GoodBuilder {
	gb.instance.Code = code
	return gb
}

func (gb GoodBuilder) Name(name string) GoodBuilder {
	gb.instance.Name = name
	return gb
}

func (gb GoodBuilder) UnitOfMeasure(uom string) GoodBuilder {
	gb.instance.UnitOfMeasure = uom
	return gb
}

func (gb GoodBuilder) Country(country string) GoodBuilder {
	gb.instance.Country = country
	return gb
}

func (gb GoodBuilder) V() models.Good {
	return *gb.instance
}

func (gb GoodBuilder) P() *models.Good {
	return gb.instance
}

func GoodsEqual(g1 *models.Good, g2 *models.Good) bool {
	if g1 == nil || g2 == nil {
		return false
	}
	return g1.Code == g2.Code && g1.Name == g2.Name && g1.Country == g2.Country && g1.UnitOfMeasure == g2.UnitOfMeasure
}
