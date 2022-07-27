package main

import (
	"context"

	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
)

func main() {
	good := goodPkg.NewWithContext(context.Background())
	go runBot((good))
	go runREST()
	runGRPCServer(good)
}
