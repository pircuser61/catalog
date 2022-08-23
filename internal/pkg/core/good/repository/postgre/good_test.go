package postgre

import (
	"context"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	goodFxtr "gitlab.ozon.dev/pircuser61/catalog/tests/fixtures"
)

func TestCreateGood(t *testing.T) {
	ctx := context.Background()

	query := "INSERT INTO good (name, unit_of_measure_id, country_id) VALUES ($1, $2, $3);"
	good := goodFxtr.Good().Code(1).Name("Good").P()
	one, id := 1, uint32(1)
	t.Run("create good success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		ExpectGetKeys(ctx, &f, good, &one, &id, &id)
		f.pool.EXPECT().Exec(ctx, query, good.Name, id, id)

		// act
		err := f.goodRepo.Add(ctx, good)

		// assert
		assert.NoError(t, err)
	})

	t.Run("create good keys error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		ExpectGetKeys(ctx, &f, good, &one, nil, &id)

		// act
		err := f.goodRepo.Add(ctx, good)

		// assert
		assert.ErrorIs(t, err, models.ErrValidation)
	})
}

func TestUpdateGood(t *testing.T) {
	ctx := context.Background()
	good := goodFxtr.Good().Code(1).Name("Good").P()
	query := "UPDATE good SET name  = $2, unit_of_measure_id = $3, country_id = $4 WHERE code = $1;"
	one, id := 1, uint32(1)
	t.Run("update good", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		ExpectGetKeys(ctx, &f, good, &one, &id, &id)
		result := pgxmock.NewResult("", 1)
		f.pool.EXPECT().Exec(ctx, query, good.Code, good.Name, id, id).Return(result, nil)

		// act
		err := f.goodRepo.Update(ctx, good)

		// assert
		require.NoError(t, err)
	})

	t.Run("update good not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		ExpectGetKeys(ctx, &f, good, &one, &id, &id)
		result := pgxmock.NewResult("", 0)
		f.pool.EXPECT().Exec(ctx, query, good.Code, good.Name, id, id).Return(result, nil)

		// act
		err := f.goodRepo.Update(ctx, good)

		// assert
		assert.ErrorIs(t, err, storePkg.ErrNotExists)
	})

	t.Run("update good keys error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		ExpectGetKeys(ctx, &f, good, &one, &id, nil)

		// act
		err := f.goodRepo.Update(ctx, good)

		// assert
		assert.ErrorIs(t, err, models.ErrValidation)
	})

	t.Run("update good error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		ExpectGetKeys(ctx, &f, good, &one, &id, &id)
		someError := errors.New("some error")

		f.pool.EXPECT().Exec(ctx, query, good.Code, good.Name, id, id).Return(nil, someError)
		// act
		err := f.goodRepo.Update(ctx, good)

		// assert
		assert.ErrorIs(t, err, someError)
	})
}

func TestGetGood(t *testing.T) {
	ctx := context.Background()
	query := "SELECT code, g.name, uom.name as unit_of_measure, c.name as country FROM good as g" +
		" left  join country as c on c.country_id  = g.country_id" +
		" left  join unit_of_measure as uom on uom.unit_of_measure_id  = g.unit_of_measure_id" +
		" WHERE code = $1;"
	code := uint64(1)
	t.Run("get good success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		good := goodFxtr.Good().Code(code).Name("Good").UnitOfMeasure("UOM").Country("Country").P()
		columns := []string{"code", "name", "unit_of_measure", "country"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(good.Code, good.Name, good.UnitOfMeasure, good.Country).ToPgxRows()
		f.pool.EXPECT().Query(ctx, query, good.Code).Return(pgxRows, nil)

		// act
		result, err := f.goodRepo.Get(ctx, good.Code)

		// assert
		require.NoError(t, err)
		assert.Equal(t, good, result)
	})

	t.Run("get good not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		pgxRows := pgxpoolmock.NewRows([]string{}).ToPgxRows()
		f.pool.EXPECT().Query(ctx, query, code).Return(pgxRows, nil)

		// act
		_, err := f.goodRepo.Get(ctx, code)

		// assert
		assert.ErrorIs(t, err, storePkg.ErrNotExists)
	})

	t.Run("get good error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		someError := errors.New("some error")
		f.pool.EXPECT().Query(ctx, query, code).Return(nil, someError)

		// act
		_, err := f.goodRepo.Get(ctx, code)

		// assert
		assert.ErrorIs(t, err, someError)
	})
}

func TestDeleteGood(t *testing.T) {
	ctx := context.Background()
	query := "DELETE FROM good WHERE code = $1;"
	code := uint64(1)
	t.Run("delete good success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		result := pgxmock.NewResult("", 1)
		f.pool.EXPECT().Exec(ctx, query, code).Return(result, nil)

		// act
		err := f.goodRepo.Delete(ctx, code)

		// assert
		require.NoError(t, err)
	})

	t.Run("delete good not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		result := pgxmock.NewResult("", 0)
		f.pool.EXPECT().Exec(ctx, query, code).Return(result, nil)

		// act
		err := f.goodRepo.Delete(ctx, code)

		// assert
		assert.ErrorIs(t, err, storePkg.ErrNotExists)
	})

	t.Run("delete error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		someError := errors.New("some error")
		f.pool.EXPECT().Exec(ctx, query, code).Return(nil, someError)

		// act
		err := f.goodRepo.Delete(ctx, code)

		// assert
		assert.ErrorIs(t, err, someError)
	})
}

func TestListGoods(t *testing.T) {
	ctx := context.Background()
	good := goodFxtr.Good().Code(1).Name("Good").P()
	goods := []*models.Good{good}
	columns := []string{"code", "name"}

	args := make([]interface{}, 0)

	t.Run("list goods", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryList := "SELECT code, good.name FROM good" +
			" LEFT OUTER JOIN country USING (country_id)" +
			" LEFT OUTER JOIN unit_of_measure USING (unit_of_measure_id)" +
			" ORDER BY good.name"
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(good.Code, good.Name).ToPgxRows()
		f.pool.EXPECT().Query(ctx, queryList, args).Return(pgxRows, nil)

		// act
		result, err := f.goodRepo.List(ctx, 0, 0)

		// assert
		require.NoError(t, err)
		assert.Equal(t, goods, result)
	})

	t.Run("list goods limit", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryList := "SELECT code, good.name FROM good" +
			" LEFT OUTER JOIN country USING (country_id)" +
			" LEFT OUTER JOIN unit_of_measure USING (unit_of_measure_id)" +
			" ORDER BY good.name LIMIT 10"
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(good.Code, good.Name).ToPgxRows()
		f.pool.EXPECT().Query(ctx, queryList, args).Return(pgxRows, nil)

		// act
		result, err := f.goodRepo.List(ctx, 10, 0)

		// assert
		require.NoError(t, err)
		assert.Equal(t, goods, result)
	})

	t.Run("list goods offset", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryList := "SELECT code, good.name FROM good" +
			" LEFT OUTER JOIN country USING (country_id)" +
			" LEFT OUTER JOIN unit_of_measure USING (unit_of_measure_id)" +
			" ORDER BY good.name OFFSET 10"
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(good.Code, good.Name).ToPgxRows()
		f.pool.EXPECT().Query(ctx, queryList, args).Return(pgxRows, nil)

		// act
		result, err := f.goodRepo.List(ctx, 0, 10)

		// assert
		require.NoError(t, err)
		assert.Equal(t, goods, result)
	})
}

func TestGetKeys(t *testing.T) {
	ctx := context.Background()
	good := goodFxtr.Good().UnitOfMeasure("UOM").Country("Country").P()
	one, id := 1, uint32(1)

	t.Run("get good keys success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		ExpectGetKeys(ctx, &f, good, &one, &id, &id)

		// act
		result, err := f.goodRepo.GetKeys(ctx, good)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &goodPkg.GoodKeys{TheOne: &one, UnitOfMeasureId: &id, CountryId: &id}, result)
	})

	t.Run("the one not exists", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		ExpectGetKeys(ctx, &f, good, nil, nil, nil)

		// act
		_, err := f.goodRepo.GetKeys(ctx, good)

		// assert
		assert.ErrorIs(t, err, storePkg.ErrNotExists)
	})

	t.Run("get keys unit not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		ExpectGetKeys(ctx, &f, good, &one, nil, &id)

		// act
		_, err := f.goodRepo.GetKeys(ctx, good)

		// assert
		assert.ErrorIs(t, err, models.ErrValidation)
		assert.Equal(t, err.Error(),
			errors.WithMessagef(models.ErrValidation,
				"Единица измерения %s не найдена в справочнике",
				good.UnitOfMeasure).Error())
	})

	t.Run("get keys country not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		ExpectGetKeys(ctx, &f, good, &one, &id, nil)

		// act
		_, err := f.goodRepo.GetKeys(ctx, good)

		// assert
		assert.ErrorIs(t, err, models.ErrValidation)
		assert.Equal(t, err.Error(),
			errors.WithMessagef(models.ErrValidation,
				"Страна %s не найдена в справочнике",
				good.Country).Error())

	})

}

func ExpectGetKeys(ctx context.Context, f *goodsRepoFixtures, good *models.Good,
	one *int, unit_id, country_id *uint32) {
	queryKeys := "SELECT 1 as the_one, (SELECT country_id FROM country WHERE country.name = $2)," +
		" (SELECT unit_of_measure_id FROM unit_of_measure WHERE unit_of_measure.name = $1);"

	columns := []string{"the_one", "country_id", "unit_of_measure_id"}
	var pgxRows pgx.Rows
	if one == nil {
		pgxRows = pgxpoolmock.NewRows([]string{}).ToPgxRows()
	} else {
		pgxRows = pgxpoolmock.NewRows(columns).AddRow(one, country_id, unit_id).ToPgxRows()
	}
	f.pool.EXPECT().Query(ctx, queryKeys, good.UnitOfMeasure, good.Country).Return(pgxRows, nil)
}
