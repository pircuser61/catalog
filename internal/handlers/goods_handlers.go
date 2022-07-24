package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"gitlab.ozon.dev/pircuser61/catalog/internal/commander"
	"gitlab.ozon.dev/pircuser61/catalog/internal/storage"
)

const (
	listCmd   = "list"
	addCmd    = "add"
	udateCmd  = "update"
	deleteCmd = "delete"
)

var BadArgument = errors.New("bad argument")

func listGoods(_ string) string {
	data := storage.List()
	if len(data) == 0 {
		return "Список пуст"
	}
	res := make([]string, len(data))
	for _, x := range data {
		res = append(res, x.String())
	}
	return strings.Join(res, "\n")
}

func addGoods(arg string) string {
	params := strings.Split(arg, " ")
	if len(params) != 3 {
		return errors.Wrapf(BadArgument, "%d items <%v>", len(params), params).Error()
	}
	g, err := storage.NewGood(params[0], params[1], params[2])
	if err != nil {
		return err.Error()
	}
	err = storage.Add(g)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Добавлен товар %s", g.GetName())
}

func updateGoods(arg string) string {
	params := strings.Split(arg, " ")
	if len(params) != 4 {
		return errors.Wrapf(BadArgument, "%d items <%v>", len(params), params).Error()
	}
	code, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return err.Error()
	}
	g, err := storage.Get(code)
	if err != nil {
		return err.Error()
	}
	err = g.SetFields(params[1], params[2], params[3])
	if err != nil {
		return err.Error()
	}
	err = storage.Update(g)
	if err != nil {
		return err.Error()
	}

	return "Изменения сохранены"
}

func deleteGoods(arg string) string {
	params := strings.Split(arg, " ")
	if len(params) != 1 {
		return errors.Wrapf(BadArgument, "%d items <%v>", len(params), params).Error()
	}
	code, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return err.Error()
	}
	err = storage.Delete(code)
	if err != nil {
		return err.Error()
	}

	return "Товар удален"
}

func Register(cmr *commander.Commander) {
	cmr.Register(listCmd, listGoods, "/list Список товаров")
	cmr.Register(addCmd, addGoods, "/add <name> <uom> <country> Добавить товар")
	cmr.Register(udateCmd, updateGoods, "/update <code> <name> <uom> <country> Изменить товар")
	cmr.Register(deleteCmd, deleteGoods, "/delete <code> Удалить товар")
}
