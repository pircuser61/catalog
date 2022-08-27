package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gitlab.ozon.dev/pircuser61/catalog/config"
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

	err := makeMigrations(ctx)
	if err != nil {
		panic("Make migration error: " + err.Error())
	}

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
	fmt.Println("DB connected")

	store := &storePkg.Core{
		Good:          goodRepo.New(pool, timeout),
		Country:       countryRepo.New(pool, timeout),
		UnitOfMeasure: unitOfMeasureRepo.New(pool, timeout),
	}

	go runBot(ctx, store.Good)
	go runKafkaConsumer(ctx, store)
	runGRPCServer(ctx, store)
}

func makeMigrations(ctx context.Context) error {
	fmt.Println("Make migrations")
	db, err := sql.Open("postgres", config.GetConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.Up(db, "./../../migrations/")
}
