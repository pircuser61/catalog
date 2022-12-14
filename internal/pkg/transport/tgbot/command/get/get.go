package get

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	commandPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/tgbot/command"
)

type command struct {
	good goodPkg.Repository
}

func New(good goodPkg.Repository) commandPkg.Interface {
	return &command{good: good}
}

func (c *command) Name() string {
	return "get"
}

func (c *command) Description() string {
	return "<code>"
}

func (c *command) Process(ctx context.Context, args string) string {
	params := strings.Split(args, " ")
	if len(params) != 1 {
		return fmt.Sprintf("invalid args %d items <%v>", len(params), params)
	}
	code, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return err.Error()
	}
	g, err := c.good.Get(ctx, code)
	if err != nil {
		if errors.Is(err, storePkg.ErrNotExists) {
			return "not found"
		}
		return "internal error"
	}
	return fmt.Sprintf("Code: %d Name: %s UOM: %s Country: %s", g.Code, g.Name, g.UnitOfMeasure, g.Country)
}
