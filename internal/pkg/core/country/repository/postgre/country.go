package postgre

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	countryPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type CountryRepository struct {
	conn    *pgxpool.Pool
	timeout time.Duration
}

func New(pgxConnextion *pgxpool.Pool, timeout time.Duration) countryPkg.Repository {
	return &CountryRepository{
		conn:    pgxConnextion,
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
	if err := pgxscan.Select(ctx, c.conn, &result, queryList); err != nil {
		return nil, fmt.Errorf("Country.List: %w", err)
	}
	return result, nil
}

func (c *CountryRepository) Add(ctx context.Context, ct *models.Country) error {
	if _, err := c.conn.Exec(ctx, queryAdd, ct.Name); err != nil {
		return fmt.Errorf("Country.Add: %w", err)
	}
	return nil
}

func (c *CountryRepository) Get(ctx context.Context, code uint32) (*models.Country, error) {
	result := models.Country{}
	if err := pgxscan.Get(ctx, c.conn, &result, queryGet, code); err != nil {
		if pgxscan.NotFound(err) {
			return nil, storePkg.ErrNotExists
		}
		return nil, fmt.Errorf("Country.Get: %w", err)
	}
	return &result, nil
}

func (c *CountryRepository) Update(ctx context.Context, ct *models.Country) error {
	if _, err := c.conn.Exec(ctx, queryUpdate, ct.CountryId, ct.Name); err != nil {
		return fmt.Errorf("Country.Update: %w", err)
	}
	return nil
}

func (c *CountryRepository) Delete(ctx context.Context, code uint32) error {
	commandTag, err := c.conn.Exec(ctx, queryDelete, code)
	if err != nil {
		return fmt.Errorf("Country.Delete: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return storePkg.ErrNotExists
	}
	return nil
}

// Метод для тестов
func (c *CountryRepository) GetByName(ctx context.Context, name string) (*models.Country, error) {
	//ctx, cancel := context.WithTimeout(ctx, c.timeout)
	//defer cancel()

	// для тестов таймаута
	slowRequest := "SELECT pg_sleep (2)"
	commandTag, err := c.conn.Exec(ctx, slowRequest)
	if err != nil {
		return nil, err
	}
	fmt.Print(commandTag.String())
	return nil, errors.New("TEST ERR")
	/*	result := models.Country{}
		fmt.Println(name) // для тестов с внедрением SQL
		if err := pgxscan.Get(ctx, c.conn, &result, queryByName, name); err != nil {
			if pgxscan.NotFound(err) {
				return nil, storePkg.ErrNotExists
			}
			return nil, fmt.Errorf("Country.Get: %w", err)
		}
		return &result, nil
	*/
}
