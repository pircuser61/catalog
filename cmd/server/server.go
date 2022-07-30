package main

import (
	"context"
	_ "embed"
	"net"
	"net/http"

	"github.com/flowchartsman/swaggerui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	apiPkg "gitlab.ozon.dev/pircuser61/catalog/internal/api"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grcpAddr    = ":8081"
	httpAddr    = ":8080"
	swaggerAddr = ":8082"
)

//go:embed swagger/api/api.swagger.json
var spec []byte

func runGRPCServer(good goodPkg.Interface) {
	listener, err := net.Listen("tcp", grcpAddr)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCatalogServer(grpcServer, apiPkg.New(good))
	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runREST(ctx context.Context) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterCatalogHandlerFromEndpoint(ctx, mux, grcpAddr, opts); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		panic(err)
	}
}

func runSwagger() {

	http.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	if err := http.ListenAndServe(swaggerAddr, nil); err != nil {
		panic(err)
	}
}
