package postgre

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

const (
	queryUnitOfMeasureList   = "SELECT unit_of_measure_id, name, FROM unit_of_measure;"
	queryUnitOfMeasureAdd    = "INSERT INTO unit_of_measure (name) VALUES ($1);"
	queryUnitOfMeasureGet    = "SELECT name FROM unit_of_measure WHERE unit_of_measure_id = $1;"
	queryUnitOfMeasureUpdate = "UPDATE unit_of_measure SET name  = $2 WHERE unit_of_measure_id = $1;"
	queryUnitOfMeasureDelete = "DELETE FROM unit_of_measure WHERE unit_of_measure_id = $1;"
)

func (c *dbPostgre) UnitOfMeasureList(ctx context.Context) ([]*models.UnitOfMeasure, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	var result []*models.UnitOfMeasure
	if err := pgxscan.Select(ctx, c.conn, &result, queryUnitOfMeasureList); err != nil {
		return nil, fmt.Errorf("UnitOfMeasure.List: select: %w", err)
	}
	return result, nil
}

func (c *dbPostgre) UnitOfMeasureAdd(ctx context.Context, ct *models.UnitOfMeasure) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if _, err := c.conn.Exec(ctx, queryUnitOfMeasureAdd, ct.Name); err != nil {
		return fmt.Errorf("UnitOfMeasure.Add: select: %w", err)
	}
	return nil
}

func (c *dbPostgre) UnitOfMeasureGet(ctx context.Context, code uint32) (*models.UnitOfMeasure, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result := models.UnitOfMeasure{}
	if err := pgxscan.Get(ctx, c.conn, &result, queryUnitOfMeasureGet, code); err != nil {
		return nil, fmt.Errorf("UnitOfMeasure.Get: %w", err)
	}
	return &result, nil
}

func (c *dbPostgre) UnitOfMeasureUpdate(ctx context.Context, ct *models.UnitOfMeasure) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if _, err := c.conn.Exec(ctx, queryUnitOfMeasureUpdate, ct.Unit_of_measure_id, ct.Name); err != nil {
		return fmt.Errorf("UnitOfMeasure.Update: select: %w", err)
	}
	return nil
}

func (c *dbPostgre) UnitOfMeasureDelete(ctx context.Context, code uint32) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if _, err := c.conn.Exec(ctx, queryUnitOfMeasureDelete, code); err != nil {
		return fmt.Errorf("UnitOfMeasure.Delete: select: %w", err)
	}
	return nil
}
