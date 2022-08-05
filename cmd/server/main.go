package main

import (
	"context"

	postgreStorePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage/postgre"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := postgreStorePkg.New(ctx)
	defer func() {
		store.Close(ctx)
	}()
	//go runBot(ctx, good)
	go runREST(ctx)
	runGRPCServer(ctx, store)
}
