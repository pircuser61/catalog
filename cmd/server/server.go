package main

import (
	"context"
	"net"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	grpcApiPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/grpc"

	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	"google.golang.org/grpc"
)

const (
	grcpAddr = ":8082"
)

func runGRPCServer(ctx context.Context, store *storePkg.Core) {
	listener, err := net.Listen("tcp", grcpAddr)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCatalogServer(grpcServer, grpcApiPkg.New(ctx, store))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
