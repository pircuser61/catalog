package postgre

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

const (
	queryList   = "SELECT code, name, unit_of_measure country FROM GOOD;"
	queryAdd    = "INSERT INTO good (name, unit_of_measure, country) VALUES ($1, $2, $3);"
	queryGet    = "SELECT code, name, unit_of_measure, country FROM good WHERE code = $1;"
	queryUpdate = "UPDATE good SET name  = $2, unit_of_measure = $3, country = $4 WHERE code = $1;"
	queryDelete = "DELETE FROM good WHERE code = $1;"
)

func (c *dbPostgre) GoodList(ctx context.Context) ([]*models.Good, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	var result []*models.Good
	if err := pgxscan.Select(ctx, c.conn, &result, queryList); err != nil {
		return nil, fmt.Errorf("Good.List: select: %w", err)
	}
	return result, nil
}

func (c *dbPostgre) GoodAdd(ctx context.Context, g *models.Good) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.conn.Exec(ctx, queryAdd, g.Name, g.UnitOfMeasure, g.Country)

	return err
}

func (c *dbPostgre) GoodGet(ctx context.Context, code uint64) (*models.Good, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result := models.Good{}
	if err := pgxscan.Get(ctx, c.conn, &result, queryGet, code); err != nil {
		return nil, fmt.Errorf("Good.Get: %w", err)
	}
	return &result, nil
}

func (c *dbPostgre) GoodUpdate(ctx context.Context, g *models.Good) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := c.conn.Exec(ctx, queryUpdate, g.Code, g.Name, g.UnitOfMeasure, g.Country)
	return err
}

func (c *dbPostgre) GoodDelete(ctx context.Context, code uint64) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := c.conn.Exec(ctx, queryDelete, code)
	return err
}
