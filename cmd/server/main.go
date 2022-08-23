package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/pircuser61/catalog/internal/config"
	countryRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country/repository/postgre"
	goodRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/repository/postgre"
	unitOfMeasureRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure/repository/postgre"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type DB struct {
	Timeout time.Duration
	Conn    *pgx.Conn
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timeout := time.Duration(time.Millisecond * 1000)
	psqlConn := config.GetConnectionString()

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		panic("Unable to connect to database: %v\n" + err.Error())
	}
	defer func() {
		fmt.Println("Disconnected")
		pool.Close()
	}()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Connected")

	store := &storePkg.Core{
		Good:          goodRepo.New(pool, timeout),
		Country:       countryRepo.New(pool, timeout),
		UnitOfMeasure: unitOfMeasureRepo.New(pool, timeout),
	}

	go runBot(ctx, store.Good)
	runGRPCServer(ctx, store)
}
