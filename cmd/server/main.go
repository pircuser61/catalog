package main

import (
	"context"
	"log"

	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//good := goodPkg.NewPostgre(ctx)
	good := goodPkg.New()
	defer func() {
		err := good.Disconnect(ctx)
		if err != nil {
			log.Panic("disconnect?")
		}
	}()
	go runBot(ctx, good)
	go runREST(ctx)
	runGRPCServer(good)
}
