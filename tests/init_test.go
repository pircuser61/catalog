//go:build integration
// +build integration

package tests

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"

	goodRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/repository/postgre"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	grpcApiPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/grpc"
	config "gitlab.ozon.dev/pircuser61/catalog/tests/config"
	"gitlab.ozon.dev/pircuser61/catalog/tests/postgres"
	"google.golang.org/grpc"
)

var (
	Db            *postgres.TestDB
	Timeout       time.Duration
	CatalogClient pb.CatalogClient
)

func init() {
	ctx := context.Background()
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	Db = postgres.NewTestDb(cfg)
	Timeout = time.Second * 2

	waitSv := make(chan struct{})
	go runTestGRPCServer(ctx, cfg, waitSv)
	<-waitSv
	conn, err := grpc.Dial(cfg.GrpcHost, grpc.WithInsecure(), grpc.WithTimeout(Timeout))
	if err != nil {
		panic(err)
	}
	CatalogClient = pb.NewCatalogClient(conn)
	fmt.Println("GRPC Client ready")
}

func runTestGRPCServer(ctx context.Context, cfg *config.Config, waitSv chan struct{}) {
	goodRepo := goodRepo.New(Db.Pool, Timeout)
	store := &storePkg.Core{
		Good:          goodRepo,
		Country:       nil,
		UnitOfMeasure: nil,
	}
	listener, err := net.Listen("tcp", cfg.GrpcHost)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCatalogServer(grpcServer, grpcApiPkg.New(ctx, store))
	fmt.Println("GRPC Server ready")
	close(waitSv)
	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
	fmt.Println("GRPC Server stoped") // как остановить при завершении тестов
}
