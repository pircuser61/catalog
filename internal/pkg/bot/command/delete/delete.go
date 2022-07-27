package delete

import (
	"errors"
	"fmt"
	"strconv"
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
	return "delete"
}

func (c *command) Description() string {
	return "<code>"
}

func (c *command) Process(args string) string {
	params := strings.Split(args, " ")
	if len(params) != 1 {
		return fmt.Sprintf("invalid args %d items <%v>", len(params), params)
	}
	code, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return err.Error()
	}

	if err := c.good.Delete(code); err != nil {
		if errors.Is(err, goodPkg.ErrNotFound) {
			return "not found"
		}
		return "internal error"
	}
	return "success"
}