package main

import (
	"net"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	apiPkg "gitlab.ozon.dev/pircuser61/catalog/internal/api"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"google.golang.org/grpc"
)

func runGRPCServer(user goodPkg.Interface) {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCatalogServer(grpcServer, apiPkg.New(user))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
