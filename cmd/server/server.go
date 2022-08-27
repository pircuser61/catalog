package main

import (
	"context"
	_ "expvar"
	"net"
	_ "net/http/pprof"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	config "gitlab.ozon.dev/pircuser61/catalog/config"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	grpcApiPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/grpc"
	"google.golang.org/grpc"
)

func runGRPCServer(ctx context.Context, store *storePkg.Core) {
	listener, err := net.Listen("tcp", config.GrpcAddr)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCatalogServer(grpcServer, grpcApiPkg.New(ctx, store))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
