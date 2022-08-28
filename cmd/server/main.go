package main

import (
	"context"
	"database/sql"
	"expvar"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/pressly/goose/v3"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/pircuser61/catalog/config"
	counters "gitlab.ozon.dev/pircuser61/catalog/internal/counters"
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

	tracer, closer, err := createTracer()
	if err != nil {
		panic(err.Error())
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	err = makeMigrations(ctx)
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
	go runHttp(ctx)
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

func runHttp(ctx context.Context) {

	counters.Init()

	expvar.Publish("Total", &counters.RequestCount)
	expvar.Publish("Success", &counters.SuccessCount)
	expvar.Publish("Errors", &counters.ErrorCount)

	if err := http.ListenAndServe(config.HttpAddr, nil); err != nil {
		panic(err)
	}
}

func createTracer() (opentracing.Tracer, io.Closer, error) {
	rc := &jaegerConfig.ReporterConfig{
		LocalAgentHostPort: config.JaegerHostPort,
		LogSpans:           true,
	}

	sc := &jaegerConfig.SamplerConfig{
		Type:  "const",
		Param: 1,
	}

	cfg := jaegerConfig.Configuration{
		ServiceName: "catalog",
		Disabled:    false,
		Reporter:    rc,
		Sampler:     sc,
	}
	return cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
}
