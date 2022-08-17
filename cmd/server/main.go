package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/pircuser61/catalog/internal/config"
	countryRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country/repository/postgre"
	countryUseCase "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country/usecase"
	goodRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/repository/postgre"
	goodUseCase "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/usecase"
	unitOfMeasureRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure/repository/postgre"
	unitOfMeasureUseCase "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure/usecase"
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

	rpGood := goodRepo.New(pool, timeout)
	rpCountry := countryRepo.New(pool, timeout)
	rpUnitUfMeasure := unitOfMeasureRepo.New(pool, timeout)
	store := &storePkg.Core{
		Good:          goodUseCase.New(rpGood, rpUnitUfMeasure, rpCountry),
		Country:       countryUseCase.New(rpCountry),
		UnitOfMeasure: unitOfMeasureUseCase.New(rpUnitUfMeasure),
	}

	go runBot(ctx, store.Good)
	runGRPCServer(ctx, store)
}
