package add

import (
	"errors"
	"fmt"
	"strings"

	commandPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

type command struct {
	good goodPkg.Interface
}

func New(good goodPkg.Interface) commandPkg.Interface {
	return &command{good: good}
}

func (c *command) Name() string {
	return "add"
}

func (c *command) Description() string {
	return "<name> <unit of measure> <country>"
}

func (c *command) Process(args string) string {
	params := strings.Split(args, " ")
	if len(params) != 3 {
		return fmt.Sprintf("invalid args %d items <%v>", len(params), params)
	}
	if err := c.good.Create(models.Good{Name: params[0], UnitOfMeasure: params[1], Country: params[2]}); err != nil {
		if errors.Is(err, goodPkg.ErrValidation) {
			return fmt.Sprintf("invalid args %d items <%v>", len(params), params)
		}
		return "internal error"
	}
	return "success"
}