package main

import (
	"context"
	"log"

	postgreStorePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage/postgre"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := postgreStorePkg.New(ctx)
	defer func() {
		err := store.Close(ctx)
		if err != nil {
			log.Panic("disconnect?")
		}
	}()
	//go runBot(ctx, good)
	go runREST(ctx)
	runGRPCServer(ctx, store)
}
