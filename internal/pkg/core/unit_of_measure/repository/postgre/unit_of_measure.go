package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	unitOfMeasurePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type UnitOfMeasureRepository struct {
	conn    *pgx.Conn
	timeout time.Duration
}

func New(pgxConnextion *pgx.Conn, timeout time.Duration) unitOfMeasurePkg.Repository {
	return &UnitOfMeasureRepository{
		conn:    pgxConnextion,
		timeout: timeout,
	}
}

const (
	queryList   = "SELECT unit_of_measure_id, name FROM unit_of_measure;"
	queryAdd    = "INSERT INTO unit_of_measure (name) VALUES ($1);"
	queryGet    = "SELECT name FROM unit_of_measure WHERE unit_of_measure_id = $1;"
	queryByName = "SELECT unit_of_measure_id, name FROM unit_of_measure WHERE name = $1;"
	queryUpdate = "UPDATE unit_of_measure SET name  = $2 WHERE unit_of_measure_id = $1;"
	queryDelete = "DELETE FROM unit_of_measure WHERE unit_of_measure_id = $1;"
)

func (c *UnitOfMeasureRepository) List(ctx context.Context) ([]*models.UnitOfMeasure, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	var result []*models.UnitOfMeasure
	if err := pgxscan.Select(ctx, c.conn, &result, queryList); err != nil {
		return nil, fmt.Errorf("UnitOfMeasure.List: select: %w", err)
	}
	return result, nil
}

func (c *UnitOfMeasureRepository) Add(ctx context.Context, ct *models.UnitOfMeasure) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if _, err := c.conn.Exec(ctx, queryAdd, ct.Name); err != nil {
		return fmt.Errorf("UnitOfMeasure.Add: select: %w", err)
	}
	return nil
}

func (c *UnitOfMeasureRepository) Get(ctx context.Context, code uint32) (*models.UnitOfMeasure, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result := models.UnitOfMeasure{}
	if err := pgxscan.Get(ctx, c.conn, &result, queryGet, code); err != nil {
		if pgxscan.NotFound(err) {
			return nil, storePkg.ErrNotExists
		}
		return nil, fmt.Errorf("UnitOfMeasure.Get: %w", err)
	}
	return &result, nil
}

func (c *UnitOfMeasureRepository) Update(ctx context.Context, ct *models.UnitOfMeasure) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if _, err := c.conn.Exec(ctx, queryUpdate, ct.UnitOfMeasureId, ct.Name); err != nil {
		return fmt.Errorf("UnitOfMeasure.Update: select: %w", err)
	}
	return nil
}

func (c *UnitOfMeasureRepository) Delete(ctx context.Context, code uint32) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	commandTag, err := c.conn.Exec(ctx, queryDelete, code)
	if err != nil {
		return fmt.Errorf("UnitOfMeasure.Delete: select: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return storePkg.ErrNotExists
	}

	return nil
}

func (c *UnitOfMeasureRepository) GetByName(ctx context.Context, name string) (*models.UnitOfMeasure, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result := models.UnitOfMeasure{}
	if err := pgxscan.Get(ctx, c.conn, &result, queryByName, name); err != nil {
		if pgxscan.NotFound(err) {
			return nil, storePkg.ErrNotExists
		}
		return nil, fmt.Errorf("UnitOfMeasure.Get: %w", err)
	}
	return &result, nil
}
