package postgre

import (
	"context"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	goodFxtr "gitlab.ozon.dev/pircuser61/catalog/tests/fixtures"
)

func TestCreateGood(t *testing.T) {
	ctx := context.Background()
	t.Run("create good", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()
		good := goodFxtr.Good().Code(1).Name("Good").P()

		var unit_of_measure_id, country_id uint32
		/*
			goodKeys := &goodPkg.GoodKeys{UnitOfMeasureId: &unit_of_measure_id, CountryId: &country_id}
		*/
		queryGet := "INSERT INTO good (name, unit_of_measure_id, country_id) VALUES ($1, $2, $3);"

		f.pool.EXPECT().Exec(ctx, queryGet, good.Name, unit_of_measure_id, country_id)

		// act
		err := f.goodRepo.Add(ctx, good)

		// assert
		require.NoError(t, err)
	})
}

func TestUpdateGood(t *testing.T) {
	ctx := context.Background()
	t.Run("update good", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()
		good := goodFxtr.Good().Code(1).Name("Good").P()
		var unit_of_measure_id, country_id uint32
		//goodKeys := &goodPkg.GoodKeys{UnitOfMeasureId: &unit_of_measure_id, CountryId: &country_id}

		queryUpdate := "UPDATE good SET name  = $2, unit_of_measure_id = $3, country_id = $4 WHERE code = $1;"
		result := pgxmock.NewResult("", 1)
		f.pool.EXPECT().Exec(ctx, queryUpdate, good.Code, good.Name, unit_of_measure_id, country_id).Return(result, nil)

		// act
		err := f.goodRepo.Update(ctx, good)

		// assert
		require.NoError(t, err)
	})
}

func TestGetGood(t *testing.T) {
	ctx := context.Background()
	t.Run("get good", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		good := goodFxtr.Good().Code(1).Name("Good").UnitOfMeasure("UOM").Country("Country").P()

		queryGet := "SELECT code, g.name, uom.name as unit_of_measure, c.name as country FROM good as g" +
			" left  join country as c on c.country_id  = g.country_id" +
			" left  join unit_of_measure as uom on uom.unit_of_measure_id  = g.unit_of_measure_id" +
			" WHERE code = $1;"

		columns := []string{"code", "name", "unit_of_measure", "country"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(good.Code, good.Name, good.UnitOfMeasure, good.Country).ToPgxRows()
		f.pool.EXPECT().Query(ctx, queryGet, good.Code).Return(pgxRows, nil)

		// act
		result, err := f.goodRepo.Get(ctx, good.Code)

		// assert
		require.NoError(t, err)
		assert.Equal(t, good, result)
	})
}

func TestListGoods(t *testing.T) {
	ctx := context.Background()
	t.Run("list goods", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		good := goodFxtr.Good().Code(1).Name("Good").P()
		goods := []*models.Good{good}

		queryList := "SELECT code, good.name FROM good" +
			" LEFT OUTER JOIN country USING (country_id)" +
			" LEFT OUTER JOIN unit_of_measure USING (unit_of_measure_id)" +
			" ORDER BY good.name"

		//queryGet = regexp.QuoteMeta(queryGet) //Ломает тест
		columns := []string{"code", "name"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(good.Code, good.Name).ToPgxRows()
		args := make([]interface{}, 0)
		f.pool.EXPECT().Query(ctx, queryList, args).Return(pgxRows, nil)

		// act
		result, err := f.goodRepo.List(ctx, 0, 0)

		// assert
		require.NoError(t, err)
		assert.Equal(t, goods, result)
	})
}

func TestGetKeys(t *testing.T) {
	ctx := context.Background()
	t.Run("get good", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		good := goodFxtr.Good().Code(1).Name("Good").UnitOfMeasure("UOM").Country("Country").P()

		queryKeys := "SELECT 1 as the_one, (SELECT country_id FROM country WHERE country.name = $2)," +
			" (SELECT unit_of_measure_id FROM unit_of_measure WHERE unit_of_measure.name = $1);"

		columns := []string{"the_one", "country_id", "unit_of_measure_id"}
		one := 1
		id := uint32(1)
		goodKeys := &goodPkg.GoodKeys{TheOne: &one, UnitOfMeasureId: &id, CountryId: &id}

		pgxRows := pgxpoolmock.NewRows(columns).AddRow(&one, &id, &id).ToPgxRows()
		f.pool.EXPECT().Query(ctx, queryKeys, good.UnitOfMeasure, good.Country).Return(pgxRows, nil)

		// act
		result, err := f.goodRepo.GetKeys(ctx, good)

		// assert
		require.NoError(t, err)
		assert.Equal(t, goodKeys, result)
	})
}
