package storage

import (
	"log"
	"strconv"

	"github.com/pkg/errors"
)

var UserNotExists = errors.New("goods does not exist")
var UserExists = errors.New("goods exist")

var data map[uint64]*Goods

func init() {
	data = make(map[uint64]*Goods)
	g, err := NewGoods("Пакет", "шт", "Россия")
	if err != nil {
		log.Panic(err)
	}
	if err := Add(g); err != nil {
		log.Panic(err)
	}

	g, err = NewGoods("Большой пакет", "шт", "Россия")
	if err != nil {
		log.Panic(err)
	}
	if err := Add(g); err != nil {
		log.Panic(err)
	}
}

func List() []*Goods {
	res := make([]*Goods, 0, len(data))
	for _, x := range data {
		res = append(res, x)
	}
	return res
}

func Add(g *Goods) error {
	if _, ok := data[g.GetCode()]; ok {
		return errors.Wrap(UserExists, strconv.FormatUint(g.GetCode(), 10))
	}
	data[g.GetCode()] = g
	return nil
}

func Get(code uint64) (*Goods, error) {
	if g, ok := data[code]; ok {
		return g, nil
	}
	return nil, errors.Wrap(UserNotExists, strconv.FormatUint(code, 10))
}

func Update(g *Goods) error {
	if _, ok := data[g.GetCode()]; !ok {
		return errors.Wrap(UserNotExists, strconv.FormatUint(g.GetCode(), 10))
	}
	data[g.GetCode()] = g
	return nil
}

func Delete(code uint64) error {
	if _, ok := data[code]; ok {
		delete(data, code)
		return nil
	}
	return errors.Wrap(UserNotExists, strconv.FormatUint(code, 10))
}
