package update

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	commandPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

type command struct {
	good goodPkg.Interface
}

func New(good goodPkg.Interface) commandPkg.Interface {
	return &command{good: good}
}

func (c *command) Name() string {
	return "update"
}

func (c *command) Description() string {
	return "<code> <name> <unit of measure> <country"
}

func (c *command) Process(ctx context.Context, args string) string {
	params := strings.Split(args, " ")
	if len(params) != 4 {
		return fmt.Sprintf("invalid args %d items <%v>", len(params), params)
	}
	code, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return err.Error()
	}
	g, err := c.good.GoodGet(ctx, code)
	if err != nil {
		if errors.Is(err, cachePkg.ErrUserNotExists) {
			return "not found"
		}
		return "internal error"
	}
	g.Name = params[1]
	g.UnitOfMeasure = params[2]
	g.Country = params[3]

	if err := c.good.GoodUpdate(ctx, *g); err != nil {
		if errors.Is(err, goodPkg.ErrNotFound) {
			return "not found"
		}
		if errors.Is(err, models.ErrValidation) {
			return err.Error()
		}
		return "internal error"
	}
	return "success"

}
