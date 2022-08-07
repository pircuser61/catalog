package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	unitOfMeasurePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type UnitOfMeasureRepository struct {
	pool    *pgxpool.Pool
	timeout time.Duration
}

func New(pgxConnextion *pgxpool.Pool, timeout time.Duration) unitOfMeasurePkg.Repository {
	return &UnitOfMeasureRepository{
		pool:    pgxConnextion,
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
	var result []*models.UnitOfMeasure
	if err := pgxscan.Select(ctx, c.pool, &result, queryList); err != nil {
		return nil, fmt.Errorf("UnitOfMeasure.List: %w", err)
	}
	return result, nil
}

func (c *UnitOfMeasureRepository) Add(ctx context.Context, ct *models.UnitOfMeasure) error {
	if _, err := c.pool.Exec(ctx, queryAdd, ct.Name); err != nil {
		return fmt.Errorf("UnitOfMeasure.Add: %w", err)
	}
	return nil
}

func (c *UnitOfMeasureRepository) Get(ctx context.Context, code uint32) (*models.UnitOfMeasure, error) {
	result := models.UnitOfMeasure{}
	if err := pgxscan.Get(ctx, c.pool, &result, queryGet, code); err != nil {
		if pgxscan.NotFound(err) {
			return nil, storePkg.ErrNotExists
		}
		return nil, fmt.Errorf("UnitOfMeasure.Get: %w", err)
	}
	return &result, nil
}

func (c *UnitOfMeasureRepository) Update(ctx context.Context, ct *models.UnitOfMeasure) error {
	if _, err := c.pool.Exec(ctx, queryUpdate, ct.UnitOfMeasureId, ct.Name); err != nil {
		return fmt.Errorf("UnitOfMeasure.Update: %w", err)
	}
	return nil
}

func (c *UnitOfMeasureRepository) Delete(ctx context.Context, code uint32) error {
	commandTag, err := c.pool.Exec(ctx, queryDelete, code)
	if err != nil {
		return fmt.Errorf("UnitOfMeasure.Delete: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return storePkg.ErrNotExists
	}

	return nil
}

func (c *UnitOfMeasureRepository) GetByName(ctx context.Context, name string) (*models.UnitOfMeasure, error) {
	result := models.UnitOfMeasure{}
	if err := pgxscan.Get(ctx, c.pool, &result, queryByName, name); err != nil {
		if pgxscan.NotFound(err) {
			return nil, storePkg.ErrNotExists
		}
		return nil, fmt.Errorf("UnitOfMeasure.Get: %w", err)
	}
	return &result, nil
}
