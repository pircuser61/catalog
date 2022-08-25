package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/georgysavva/scany/pgxscan"
	countryPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type CountryRepository struct {
	pool    pgxpoolmock.PgxPool
	timeout time.Duration
}

func New(pgxConnextion pgxpoolmock.PgxPool, timeout time.Duration) countryPkg.Repository {
	return &CountryRepository{
		pool:    pgxConnextion,
		timeout: timeout,
	}
}

const (
	queryList   = "SELECT country_id, name FROM country;"
	queryAdd    = "INSERT INTO country (name) VALUES ($1);"
	queryGet    = "SELECT name FROM country WHERE country_id = $1;"
	queryByName = "SELECT country_id, name FROM country WHERE name = $1;"
	queryUpdate = "UPDATE country SET name  = $2 WHERE country_id = $1;"
	queryDelete = "DELETE FROM country WHERE country_id = $1;"
)

func (c *CountryRepository) List(ctx context.Context) ([]*models.Country, error) {
	var result []*models.Country
	if err := pgxscan.Select(ctx, c.pool, &result, queryList); err != nil {
		return nil, fmt.Errorf("Country.List: %w", err)
	}
	return result, nil
}

func (c *CountryRepository) Add(ctx context.Context, ct *models.Country) error {
	if _, err := c.pool.Exec(ctx, queryAdd, ct.Name); err != nil {
		return fmt.Errorf("Country.Add: %w", err)
	}
	return nil
}

func (c *CountryRepository) Get(ctx context.Context, code uint32) (*models.Country, error) {
	result := models.Country{}
	if err := pgxscan.Get(ctx, c.pool, &result, queryGet, code); err != nil {
		if pgxscan.NotFound(err) {
			return nil, storePkg.ErrNotExists
		}
		return nil, fmt.Errorf("Country.Get: %w", err)
	}
	return &result, nil
}

func (c *CountryRepository) Update(ctx context.Context, ct *models.Country) error {
	commandTag, err := c.pool.Exec(ctx, queryUpdate, ct.CountryId, ct.Name)
	if err != nil {
		return fmt.Errorf("Country.Update: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return storePkg.ErrNotExists
	}
	return nil
}

func (c *CountryRepository) Delete(ctx context.Context, code uint32) error {
	commandTag, err := c.pool.Exec(ctx, queryDelete, code)
	if err != nil {
		return fmt.Errorf("Country.Delete: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return storePkg.ErrNotExists
	}
	return nil
}
