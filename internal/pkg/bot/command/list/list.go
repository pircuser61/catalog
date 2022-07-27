package GoodListResponse

import (
	"strings"

	commandPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
)

type command struct {
	good goodPkg.Interface
}

func New(good goodPkg.Interface) commandPkg.Interface {
	return &command{good: good}
}

func (c *command) Name() string {
	return "list"
}

func (c *command) Description() string {
	return "no params"
}

func (c *command) Process(args string) string {
	data := c.good.List()
	if len(data) == 0 {
		return "Список пуст"
	}
	res := make([]string, len(data))
	for _, x := range data {
		res = append(res, x.String())
	}
	return strings.Join(res, "\n")
}
