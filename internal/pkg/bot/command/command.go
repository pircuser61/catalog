package command

import "context"

type Interface interface {
	Name() string
	Description() string
	Process(context.Context, string) string
}
