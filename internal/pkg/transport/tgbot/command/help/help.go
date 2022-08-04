package help

import (
	"context"
	"fmt"

	commandPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/tgbot/command"
)

type command struct {
	extended map[string]string
}

func New(extendedMap map[string]string) commandPkg.Interface {
	if extendedMap == nil {
		extendedMap = map[string]string{}
	}
	return &command{extended: extendedMap}
}

func (c *command) Name() string {
	return "help"
}

func (c *command) Description() string {
	return "list of commands"
}

func (c *command) Process(_ context.Context, _ string) string {
	result := fmt.Sprintf("/%s - %s", c.Name(), c.Description())
	for cmd, descr := range c.extended {
		result += fmt.Sprintf("\n/%s - %s", cmd, descr)
	}

	return result
}
