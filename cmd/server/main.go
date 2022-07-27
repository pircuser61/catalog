package main

import (
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
)

func main() {
	good := goodPkg.New()
	go runBot((good))
	runGRPCServer(good)
}
