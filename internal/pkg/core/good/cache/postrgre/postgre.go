package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/pircuser61/catalog/internal/config"
	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

const (
	queryList   = "SELECT code, name, unit_of_measure country FROM GOOD;"
	queryAdd    = "INSERT INTO good (name, unit_of_measure, country) VALUES ($1, $2, $3);"
	queryGet    = "SELECT code, name, unit_of_measure, country FROM good WHERE code = $1;"
	queryUpdate = "UPDATE good SET name  = $2, unit_of_measure = $3, country = $4 WHERE code = $1;"
	queryDelete = "DELETE FROM good WHERE code = $1;"
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

func (c *dbPostgre) List(ctx context.Context) ([]*models.Good, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	var result []*models.Good
	if err := pgxscan.Select(ctx, c.conn, &result, queryList); err != nil {
		return nil, fmt.Errorf("Good.List: select: %w", err)
	}
	defer cancel()
	return result, nil
}

func (c *dbPostgre) Add(ctx context.Context, g *models.Good) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.conn.Exec(ctx, queryAdd, g.Name, g.UnitOfMeasure, g.Country)

	return err
}

func (c *dbPostgre) Get(ctx context.Context, code uint64) (*models.Good, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result := models.Good{}
	if err := pgxscan.Get(ctx, c.conn, &result, queryGet, code); err != nil {
		return nil, fmt.Errorf("Good.Get: %w", err)
	}
	return &result, nil
}

func (c *dbPostgre) Update(ctx context.Context, g *models.Good) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := c.conn.Exec(ctx, queryUpdate, g.Code, g.Name, g.UnitOfMeasure, g.Country)
	return err
}

func (c *dbPostgre) Delete(ctx context.Context, code uint64) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := c.conn.Exec(ctx, queryDelete, code)
	return err
}

func (c *dbPostgre) Disconnect(ctx context.Context) error {
	defer fmt.Println("Disconnected")
	return c.conn.Close(ctx)
}

func (c *dbPostgre) Lock() string {
	return ""
}

func (c *dbPostgre) RLock() string {
	return ""
}

func (c *dbPostgre) Unlock() (result string) {
	return ""
}

func (c *dbPostgre) RUnlock() (result string) {
	return ""
}
