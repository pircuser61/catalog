package main

import (
	"context"

	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	stopCh := make(chan struct{})

	AddContextInterceptor := func(grcpCtx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(AddContextInterceptor))

	good := goodPkg.New()

	go runBot(ctx, stopCh, good)
	go runREST(ctx)
	go runGRPCServer(ctx, grpcServer, good)
	defer func() { cancel(); grpcServer.GracefulStop() }()
	<-stopCh
}
