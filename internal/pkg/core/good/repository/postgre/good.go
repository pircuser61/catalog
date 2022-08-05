package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
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
	queryList    = queryJoin + ";"
	queryAdd     = "INSERT INTO good (name, unit_of_measure_id, country_id) VALUES ($1, $2, $3);"
	queryGet     = queryJoin + "\nWHERE code = $1;"
	queryUpdate  = "UPDATE good SET name  = $2, unit_of_measure_id = $3, country_id = $4 WHERE code = $1;"
	queryDelete  = "DELETE FROM good WHERE code = $1;"
	queryGetKeys = "SELECT 1 as the_one, (SELECT country_id FROM country WHERE country.name = $2), " +
		" (SELECT unit_of_measure_id FROM unit_of_measure WHERE unit_of_measure.name = $1);"
)

func (c *GoodsRepository) Add(ctx context.Context, good *models.Good, keys *goodPkg.GoodKeys) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := c.conn.Exec(ctx, queryAdd, good.Name, *keys.UnitOfMeasureId, *keys.CountryId)
	return err
}

func (c *GoodsRepository) Update(ctx context.Context, good *models.Good, keys *goodPkg.GoodKeys) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	commandTag, err := c.conn.Exec(ctx, queryUpdate, good.Code, good.Name, keys.UnitOfMeasureId, keys.CountryId)
	if err != nil {
		return fmt.Errorf("Good.Update: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return storePkg.ErrNotExists
	}
	return nil
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
	commandTag, err := c.conn.Exec(ctx, queryDelete, code)
	if err != nil {
		return fmt.Errorf("Good.Delete: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return storePkg.ErrNotExists
	}
	return nil
}

func (c *GoodsRepository) List(ctx context.Context, limit uint64, offset uint64) ([]*models.Good, error) {
	qBuilder := squirrel.Select("code, good.name").
		From("good").
		JoinClause("LEFT OUTER JOIN country USING (country_id)").
		JoinClause("LEFT OUTER JOIN unit_of_measure USING (unit_of_measure_id)").
		OrderBy("good.name")

	if limit > 0 {
		qBuilder = qBuilder.Limit(limit)
	}
	if offset > 0 {
		qBuilder = qBuilder.Offset(offset)
	}

	query, args, err := qBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	var result []*models.Good
	if err := pgxscan.Select(ctx, c.conn, &result, query, args...); err != nil {
		return nil, fmt.Errorf("Good.List:  %w", err)
	}
	return result, nil
}

func (c *GoodsRepository) GetKeys(ctx context.Context, good *models.Good) (*goodPkg.GoodKeys, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	result := goodPkg.GoodKeys{}

	if err := pgxscan.Get(ctx, c.conn, &result, queryGetKeys, good.UnitOfMeasure, good.Country); err != nil {
		if pgxscan.NotFound(err) {
			return nil, storePkg.ErrNotExists
		}
		return nil, fmt.Errorf("Good.GetKeys: %w", err)
	}

	if result.UnitOfMeasureId == nil {
		return nil, errors.WithMessagef(models.ErrValidation,
			"Единица измерения %s не найдена в справочнике",
			good.UnitOfMeasure)
	}
	if result.CountryId == nil {
		return nil, errors.WithMessagef(models.ErrValidation,
			"Страна %s не найдена в справочнике",
			good.Country)
	}

	return &result, nil
}
