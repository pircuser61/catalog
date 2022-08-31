package main

import (
	"context"
	"database/sql"
	"expvar"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/pressly/goose/v3"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/pircuser61/catalog/config"
	counters "gitlab.ozon.dev/pircuser61/catalog/internal/counters"
	logger "gitlab.ozon.dev/pircuser61/catalog/internal/logger"
	countryRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country/repository/postgre"
	goodRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/repository/postgre"
	unitOfMeasureRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure/repository/postgre"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	"go.uber.org/zap"
)

type DB struct {
	Timeout time.Duration
	Conn    *pgx.Conn
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer logger.Sync()

	tracer, closer, err := createTracer()
	if err != nil {
		logger.Panic("Create tracer error:" + err.Error())
	}
	defer closer.Close()

	// NewLoggingTracer creates a new tracer that logs all span interactions
	opentracing.SetGlobalTracer(tracer)

	err = makeMigrations(ctx)
	if err != nil {
		logger.Panic("Make migration error:" + err.Error())
	}

	timeout := time.Duration(time.Millisecond * 1000)
	logger.Debug("DB connect")
	psqlConn := config.GetConnectionString()
	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		logger.Panic("Unable to connect to database", zap.Error(err))
	}
	defer func() {
		pool.Close()
	}()

	if err := pool.Ping(ctx); err != nil {
		logger.Panic("Db ping", zap.Error(err))
	}

	store := &storePkg.Core{
		Good:          goodRepo.New(pool, timeout),
		Country:       countryRepo.New(pool, timeout),
		UnitOfMeasure: unitOfMeasureRepo.New(pool, timeout),
	}
	logger.Debug("Runing services")
	go runBot(ctx, store.Good)
	go runKafkaConsumer(ctx, store)
	go runHttp(ctx)
	runGRPCServer(ctx, store)
}

func makeMigrations(ctx context.Context) error {
	defer logger.Sync()
	logger.Debug("Make migreations")
	db, err := sql.Open("postgres", config.GetConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.Up(db, "./../../migrations/")
}

func runHttp(ctx context.Context) {
	defer logger.Sync()
	logger.Debug("Run http server")
	counters.Init()
	expvar.Publish("Total", &counters.RequestCount)
	expvar.Publish("Success", &counters.SuccessCount)
	expvar.Publish("Errors", &counters.ErrorCount)

	if err := http.ListenAndServe(config.HttpAddr, nil); err != nil {
		logger.Panic("Http listen", zap.Error(err))
	}
}

func createTracer() (opentracing.Tracer, io.Closer, error) {
	defer logger.Sync()
	logger.Debug("Create tracer")
	rc := &jaegerConfig.ReporterConfig{
		LocalAgentHostPort: config.JaegerHostPort,
		LogSpans:           false,
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
	return cfg.NewTracer(jaegerConfig.Logger(&logger.JgLogger))
}
