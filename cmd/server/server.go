package main

import (
	"context"
	_ "expvar"
	"net"
	_ "net/http/pprof"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	config "gitlab.ozon.dev/pircuser61/catalog/config"
	counters "gitlab.ozon.dev/pircuser61/catalog/internal/counters"
	logger "gitlab.ozon.dev/pircuser61/catalog/internal/logger"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	grpcApiPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func runGRPCServer(ctx context.Context, store *storePkg.Core) {
	listener, err := net.Listen("tcp", config.GrpcAddr)
	if err != nil {
		logger.Panic("", zap.Error(err))
	}

	interceptors := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
		counterInterceptor))

	grpcServer := grpc.NewServer(interceptors)
	pb.RegisterCatalogServer(grpcServer, grpcApiPkg.New(ctx, store))

	if err = grpcServer.Serve(listener); err != nil {
		logger.Panic("GRPC.Serve", zap.Error(err))
	}
}

func counterInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	counters.Request()
	result, err := handler(ctx, req)
	if err == nil {
		counters.Success()
	} else {
		counters.Error()
	}
	return result, err
}
