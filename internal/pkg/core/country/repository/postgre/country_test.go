package postgre

import (
	"context"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/pashagolub/pgxmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

func TestCreateCountry(t *testing.T) {
	ctx := context.Background()
	query := "INSERT INTO country (name) VALUES ($1);"
	country := &models.Country{Name: "Country"}
	t.Run("create country success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		f.pool.EXPECT().Exec(ctx, query, country.Name)

		// act
		err := f.countryRepo.Add(ctx, country)

		// assert
		assert.NoError(t, err)
	})
}

func TestUpdateCountry(t *testing.T) {
	ctx := context.Background()
	country := &models.Country{Name: "Country"}
	query := "UPDATE country SET name  = $2 WHERE country_id = $1;"
	t.Run("update country", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		result := pgxmock.NewResult("", 1)
		f.pool.EXPECT().Exec(ctx, query, country.CountryId, country.Name).Return(result, nil)

		// act
		err := f.countryRepo.Update(ctx, country)

		// assert
		require.NoError(t, err)
	})

	t.Run("update country not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		result := pgxmock.NewResult("", 0)
		f.pool.EXPECT().Exec(ctx, query, country.CountryId, country.Name).Return(result, nil)

		// act
		err := f.countryRepo.Update(ctx, country)

		// assert
		assert.ErrorIs(t, err, storePkg.ErrNotExists)
	})

	t.Run("update country error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		someError := errors.New("some error")

		f.pool.EXPECT().Exec(ctx, query, country.CountryId, country.Name).Return(nil, someError)
		// act
		err := f.countryRepo.Update(ctx, country)

		// assert
		assert.ErrorIs(t, err, someError)
	})
}

func TestGetCountry(t *testing.T) {
	ctx := context.Background()
	query := "SELECT name FROM country WHERE country_id = $1;"
	id := uint32(1)
	t.Run("get country success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		country := &models.Country{Name: "Country"}
		columns := []string{"country_id", "name"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(country.CountryId, country.Name).ToPgxRows()
		f.pool.EXPECT().Query(ctx, query, country.CountryId).Return(pgxRows, nil)

		// act
		result, err := f.countryRepo.Get(ctx, country.CountryId)

		// assert
		require.NoError(t, err)
		assert.Equal(t, country, result)
	})

	t.Run("get country not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		pgxRows := pgxpoolmock.NewRows([]string{}).ToPgxRows()
		f.pool.EXPECT().Query(ctx, query, id).Return(pgxRows, nil)

		// act
		_, err := f.countryRepo.Get(ctx, id)

		// assert
		assert.ErrorIs(t, err, storePkg.ErrNotExists)
	})

	t.Run("get country error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		someError := errors.New("some error")
		f.pool.EXPECT().Query(ctx, query, id).Return(nil, someError)

		// act
		_, err := f.countryRepo.Get(ctx, id)

		// assert
		assert.ErrorIs(t, err, someError)
	})
}

func TestDeleteCountry(t *testing.T) {
	ctx := context.Background()
	query := "DELETE FROM country WHERE country_id = $1;"
	country_id := uint32(1)
	t.Run("delete country success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		result := pgxmock.NewResult("", 1)
		f.pool.EXPECT().Exec(ctx, query, country_id).Return(result, nil)

		// act
		err := f.countryRepo.Delete(ctx, country_id)

		// assert
		require.NoError(t, err)
	})

	t.Run("delete country not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		result := pgxmock.NewResult("", 0)
		f.pool.EXPECT().Exec(ctx, query, country_id).Return(result, nil)

		// act
		err := f.countryRepo.Delete(ctx, country_id)

		// assert
		assert.ErrorIs(t, err, storePkg.ErrNotExists)
	})

	t.Run("delete error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		someError := errors.New("some error")
		f.pool.EXPECT().Exec(ctx, query, country_id).Return(nil, someError)

		// act
		err := f.countryRepo.Delete(ctx, country_id)

		// assert
		assert.ErrorIs(t, err, someError)
	})
}

func TestListCountrys(t *testing.T) {
	ctx := context.Background()
	country := &models.Country{Name: "Country"}
	countrys := []*models.Country{country}
	columns := []string{"country_id", "name"}
	args := make([]interface{}, 0)

	t.Run("list countrys", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryList := "SELECT country_id, name FROM country;"
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(country.CountryId, country.Name).ToPgxRows()
		f.pool.EXPECT().Query(ctx, queryList, args).Return(pgxRows, nil)

		// act
		result, err := f.countryRepo.List(ctx)

		// assert
		require.NoError(t, err)
		assert.Equal(t, countrys, result)
	})
}
