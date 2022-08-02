package postgre

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

const (
	queryCountryList   = "SELECT country_id, name, FROM country;"
	queryCountryAdd    = "INSERT INTO country (name) VALUES ($1);"
	queryCountryGet    = "SELECT name FROM country WHERE country_id = $1;"
	queryCountryUpdate = "UPDATE country SET name  = $2 WHERE country_id = $1;"
	queryCountryDelete = "DELETE FROM country WHERE country_id = $1;"
)

func (c *dbPostgre) CountryList(ctx context.Context) ([]*models.Country, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	var result []*models.Country
	if err := pgxscan.Select(ctx, c.conn, &result, queryCountryList); err != nil {
		return nil, fmt.Errorf("Country.List: select: %w", err)
	}
	return result, nil
}

func (c *dbPostgre) CountryAdd(ctx context.Context, ct *models.Country) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if _, err := c.conn.Exec(ctx, queryCountryAdd, ct.Name); err != nil {
		return fmt.Errorf("Country.Add: select: %w", err)
	}
	return nil
}

func (c *dbPostgre) CountryGet(ctx context.Context, code uint32) (*models.Country, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result := models.Country{}
	if err := pgxscan.Get(ctx, c.conn, &result, queryCountryGet, code); err != nil {
		return nil, fmt.Errorf("Country.Get: %w", err)
	}
	return &result, nil
}

func (c *dbPostgre) CountryUpdate(ctx context.Context, ct *models.Country) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if _, err := c.conn.Exec(ctx, queryCountryUpdate, ct.Country_id, ct.Name); err != nil {
		return fmt.Errorf("Country.Update: select: %w", err)
	}
	return nil
}

func (c *dbPostgre) CountryDelete(ctx context.Context, code uint32) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if _, err := c.conn.Exec(ctx, queryCountryDelete, code); err != nil {
		return fmt.Errorf("Country.Delete: select: %w", err)
	}
	return nil
}
