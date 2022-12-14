package list

import (
	"context"
	"strings"

	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	commandPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/tgbot/command"
)

type command struct {
	good goodPkg.Repository
}

func New(good goodPkg.Repository) commandPkg.Interface {
	return &command{good: good}
}

func (c *command) Name() string {
	return "list"
}

func (c *command) Description() string {
	return "no params"
}

func (c *command) Process(ctx context.Context, args string) string {
	data, err := c.good.List(ctx, 0, 0)
	if err != nil {
		return err.Error()
	}
	if len(data) == 0 {
		return "Список пуст"
	}
	res := make([]string, len(data))
	for _, x := range data {
		res = append(res, x.Name)
	}
	return strings.Join(res, "\n")
}
