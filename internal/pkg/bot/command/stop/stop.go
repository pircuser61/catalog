package stop

// Для тестов

import (
	"context"

	commandPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command"
)

type command struct {
	stopCh chan<- struct{}
}

func New(stopCh chan<- struct{}) commandPkg.Interface {
	return &command{stopCh: stopCh}
}

func (c *command) Name() string {
	return "stop"
}

func (c *command) Description() string {
	return "stops GRPC server"
}

func (c *command) Process(ctx context.Context, args string) string {
	c.stopCh <- struct{}{}
	return "SEND STOP"
}
