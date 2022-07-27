package lock

import (
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
	return "lock"
}

func (c *command) Description() string {
	return "avail commands: Rlock, Wlock, RUnlock, WUnlock"
}

func (c *command) Process(args string) string {
	switch args {
	case "Rlock":
		return c.good.GetCache().RLock()
	case "Wlock":
		return c.good.GetCache().Lock()
	case "RUnlock":
		return c.good.GetCache().RUnlock()
	case "WUnlock":
		return c.good.GetCache().Unlock()
	}
	return "wrong command " + args
}
