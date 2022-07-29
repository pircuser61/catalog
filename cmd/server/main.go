package main

import (
	"context"

	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	good := goodPkg.New()
	go runBot(ctx, good)
	go runREST(ctx)
	runGRPCServer(good)
}
