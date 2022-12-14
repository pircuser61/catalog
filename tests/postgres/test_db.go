//go:build integration
// +build integration

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	config "gitlab.ozon.dev/pircuser61/catalog/tests/config"
)

type TestDB struct {
	sync.Mutex
	Pool *pgxpool.Pool
}

type UnitOfMeasure struct {
	Id   uint32
	Name string
}
type Country struct {
	Id   uint32
	Name string
}

func NewTestDb(cfg *config.Config) *TestDB {
	ctx := context.Background()

	fmt.Println("Make migrations", cfg.DdConnectionString())
	err := makeMigrations(ctx, cfg)
	if err != nil {
		panic("Make migration error: " + err.Error())
	}
	fmt.Println("Migrations OK")

	pool, err := pgxpool.Connect(ctx, cfg.DdConnectionString())
	if err != nil {
		panic("Unable to connect to database: %v\n" + err.Error())
	}
	fmt.Println("Connected to test DB")
	return &TestDB{Pool: pool}
}

func (d *TestDB) SetUp(ctx context.Context, t *testing.T) {
	d.Lock()
	d.Truncate(ctx, "good")
}

func (d *TestDB) TearDown(ctx context.Context) {
	defer d.Unlock()
	d.Truncate(ctx, "good")
}

func (d *TestDB) Truncate(ctx context.Context, listTables string) {
	q := fmt.Sprintf("Truncate table %s", listTables)
	if _, err := d.Pool.Exec(ctx, q); err != nil {
		panic(err)
	}
}

func makeMigrations(ctx context.Context, cfg *config.Config) error {
	const migrations = "./../migrations/"

	db, err := sql.Open("postgres", cfg.DdConnectionString())
	if err != nil {
		return err
	}
	fmt.Println("Connected to ", cfg.DbName)
	defer db.Close()
	_ = goose.Down(db, migrations) /* на пустой базе дает ошибку */
	return goose.Up(db, migrations)
}

func (d *TestDB) GetFirstUnitOfMeasure(ctx context.Context, t *testing.T) *UnitOfMeasure {
	var unit_of_measure UnitOfMeasure
	if err := pgxscan.Get(ctx, d.Pool, &unit_of_measure, "SELECT unit_of_measure_id as id, name FROM unit_of_measure LIMIT 1"); err != nil {
		panic("Cant get uom id: " + err.Error())
	}
	return &unit_of_measure
}

func (d *TestDB) GetFirstCountry(ctx context.Context, t *testing.T) *Country {
	var country Country
	if err := pgxscan.Get(ctx, d.Pool, &country, "SELECT country_id as id, name FROM country LIMIT 1"); err != nil {
		panic("Cant get country id: " + err.Error())
	}
	return &country
}
