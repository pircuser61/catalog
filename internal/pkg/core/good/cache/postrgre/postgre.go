package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/pircuser61/catalog/internal/config"
	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
)

type dbPostgre struct {
	timeout time.Duration
	conn    *pgx.Conn
}

func New(ctx context.Context) cachePkg.Interface {
	tm := time.Duration(time.Millisecond * 8000)
	psqlConn := config.GetConnectionString()
	conn, err := pgx.Connect(ctx, psqlConn)
	if err != nil {
		panic("Unable to connect to database: %v\n" + err.Error())
	}
	if err := conn.Ping(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Connected")
	return &dbPostgre{timeout: tm, conn: conn}
}

func (c *dbPostgre) Close(ctx context.Context) error {
	defer fmt.Println("Disconnected")
	return c.conn.Close(ctx)
}
