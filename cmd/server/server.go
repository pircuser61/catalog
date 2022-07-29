package main

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	apiPkg "gitlab.ozon.dev/pircuser61/catalog/internal/api"

	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runGRPCServer(ctx context.Context, grpcServer *grpc.Server, good goodPkg.Interface) {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	pb.RegisterCatalogServer(grpcServer, apiPkg.New(good))
	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runREST(ctx context.Context) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterCatalogHandlerFromEndpoint(ctx, mux, ":8081", opts); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
