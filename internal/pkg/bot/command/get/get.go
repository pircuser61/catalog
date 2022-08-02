package get

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	commandPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
)

type command struct {
	good goodPkg.Interface
}

func New(good goodPkg.Interface) commandPkg.Interface {
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
	g, err := c.good.GoodGet(ctx, code)
	if err != nil {
		if errors.Is(err, cachePkg.ErrObjNotExists) {
			return "not found"
		}
		return "internal error"
	}
	return fmt.Sprintf("Code: %d Name: %s UOM: %s Country: %s", g.Code, g.Name, g.UnitOfMeasure, g.Country)
}
