package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type GoodsRepository struct {
	conn    *pgx.Conn
	timeout time.Duration
}

func New(pgxConnextion *pgx.Conn, timeout time.Duration) goodPkg.Repository {
	return &GoodsRepository{
		conn:    pgxConnextion,
		timeout: timeout,
	}
}

const (
	queryJoin = "SELECT code, g.name, uom.name as unit_of_measure, c.name as country FROM good as g" +
		"\nleft  join country as c on c.country_id  = g.country_id" +
		"\nleft  join unit_of_measure as uom on uom.unit_of_measure_id  = g.unit_of_measure_id"
	queryList   = queryJoin + ";"
	queryAdd    = "INSERT INTO good (name, unit_of_measure_id, country_id) VALUES ($1, $2, $3);"
	queryGet    = queryJoin + "\nWHERE code = $1;"
	queryUpdate = "UPDATE good SET name  = $2, unit_of_measure_id = $3, country_id = $4 WHERE code = $1;"
	queryDelete = "DELETE FROM good WHERE code = $1;"
)

func (c *GoodsRepository) Add(ctx context.Context, name string, uom_id uint32, country_id uint32) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := c.conn.Exec(ctx, queryAdd, name, uom_id, country_id)

	return err
}

func (c *GoodsRepository) Update(ctx context.Context, code uint64, name string, uom_id uint32, country_id uint32) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := c.conn.Exec(ctx, queryUpdate, code, name, uom_id, country_id)
	return err
}

func (c *GoodsRepository) Get(ctx context.Context, code uint64) (*models.Good, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result := models.Good{}

	if err := pgxscan.Get(ctx, c.conn, &result, queryGet, code); err != nil {
		if pgxscan.NotFound(err) {
			return nil, storePkg.ErrNotExists
		}
		return nil, fmt.Errorf("Good.Get: %w", err)
	}
	return &result, nil
}

func (c *GoodsRepository) Delete(ctx context.Context, code uint64) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := c.conn.Exec(ctx, queryDelete, code)
	return err
}

func (c *GoodsRepository) List(ctx context.Context) ([]*models.Good, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	var result []*models.Good
	if err := pgxscan.Select(ctx, c.conn, &result, queryList); err != nil {
		return nil, fmt.Errorf("Good.List: select: %w", err)
	}
	return result, nil
}
