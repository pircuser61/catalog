package storage

import (
	"strconv"

	"github.com/pkg/errors"
)

var UserNotExists = errors.New("good does not exist")
var UserExists = errors.New("good exist")

var data map[uint64]*Good

func init() {
	data = make(map[uint64]*Good)
}

func List() []*Good {
	res := make([]*Good, 0, len(data))
	for _, x := range data {
		res = append(res, x)
	}
	return res
}

func Add(g *Good) error {
	if _, ok := data[g.GetCode()]; ok {
		return errors.Wrap(UserExists, strconv.FormatUint(g.GetCode(), 10))
	}
	data[g.GetCode()] = g
	return nil
}

func Get(code uint64) (*Good, error) {
	if g, ok := data[code]; ok {
		return g, nil
	}
	return nil, errors.Wrap(UserNotExists, strconv.FormatUint(code, 10))
}

func Update(g *Good) error {
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
