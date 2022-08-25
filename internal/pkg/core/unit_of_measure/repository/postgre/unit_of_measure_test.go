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

func TestCreateUnitOfMeasure(t *testing.T) {
	ctx := context.Background()
	query := "INSERT INTO unit_of_measure (name) VALUES ($1);"
	unit_of_measure := &models.UnitOfMeasure{Name: "UnitOfMeasure"}
	t.Run("create unit_of_measure success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		f.pool.EXPECT().Exec(ctx, query, unit_of_measure.Name)

		// act
		err := f.unitRepo.Add(ctx, unit_of_measure)

		// assert
		assert.NoError(t, err)
	})
}

func TestUpdateUnitOfMeasure(t *testing.T) {
	ctx := context.Background()
	unit_of_measure := &models.UnitOfMeasure{Name: "UnitOfMeasure"}
	query := "UPDATE unit_of_measure SET name  = $2 WHERE unit_of_measure_id = $1;"
	t.Run("update unit_of_measure", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		result := pgxmock.NewResult("", 1)
		f.pool.EXPECT().Exec(ctx, query, unit_of_measure.UnitOfMeasureId, unit_of_measure.Name).Return(result, nil)

		// act
		err := f.unitRepo.Update(ctx, unit_of_measure)

		// assert
		require.NoError(t, err)
	})

	t.Run("update unit_of_measure not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		result := pgxmock.NewResult("", 0)
		f.pool.EXPECT().Exec(ctx, query, unit_of_measure.UnitOfMeasureId, unit_of_measure.Name).Return(result, nil)

		// act
		err := f.unitRepo.Update(ctx, unit_of_measure)

		// assert
		assert.ErrorIs(t, err, storePkg.ErrNotExists)
	})

	t.Run("update unit_of_measure error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		someError := errors.New("some error")

		f.pool.EXPECT().Exec(ctx, query, unit_of_measure.UnitOfMeasureId, unit_of_measure.Name).Return(nil, someError)
		// act
		err := f.unitRepo.Update(ctx, unit_of_measure)

		// assert
		assert.ErrorIs(t, err, someError)
	})
}

func TestGetUnitOfMeasure(t *testing.T) {
	ctx := context.Background()
	query := "SELECT name FROM unit_of_measure WHERE unit_of_measure_id = $1;"
	id := uint32(1)
	t.Run("get unit_of_measure success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		unit_of_measure := &models.UnitOfMeasure{Name: "UnitOfMeasure"}
		columns := []string{"unit_of_measure_id", "name"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(unit_of_measure.UnitOfMeasureId, unit_of_measure.Name).ToPgxRows()
		f.pool.EXPECT().Query(ctx, query, unit_of_measure.UnitOfMeasureId).Return(pgxRows, nil)

		// act
		result, err := f.unitRepo.Get(ctx, unit_of_measure.UnitOfMeasureId)

		// assert
		require.NoError(t, err)
		assert.Equal(t, unit_of_measure, result)
	})

	t.Run("get unit_of_measure not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		pgxRows := pgxpoolmock.NewRows([]string{}).ToPgxRows()
		f.pool.EXPECT().Query(ctx, query, id).Return(pgxRows, nil)

		// act
		_, err := f.unitRepo.Get(ctx, id)

		// assert
		assert.ErrorIs(t, err, storePkg.ErrNotExists)
	})

	t.Run("get unit_of_measure error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		someError := errors.New("some error")
		f.pool.EXPECT().Query(ctx, query, id).Return(nil, someError)

		// act
		_, err := f.unitRepo.Get(ctx, id)

		// assert
		assert.ErrorIs(t, err, someError)
	})
}

func TestDeleteUnitOfMeasure(t *testing.T) {
	ctx := context.Background()
	query := "DELETE FROM unit_of_measure WHERE unit_of_measure_id = $1;"
	unit_of_measure_id := uint32(1)
	t.Run("delete unit_of_measure success", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		result := pgxmock.NewResult("", 1)
		f.pool.EXPECT().Exec(ctx, query, unit_of_measure_id).Return(result, nil)

		// act
		err := f.unitRepo.Delete(ctx, unit_of_measure_id)

		// assert
		require.NoError(t, err)
	})

	t.Run("delete unit_of_measure not found", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		result := pgxmock.NewResult("", 0)
		f.pool.EXPECT().Exec(ctx, query, unit_of_measure_id).Return(result, nil)

		// act
		err := f.unitRepo.Delete(ctx, unit_of_measure_id)

		// assert
		assert.ErrorIs(t, err, storePkg.ErrNotExists)
	})

	t.Run("delete error", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		someError := errors.New("some error")
		f.pool.EXPECT().Exec(ctx, query, unit_of_measure_id).Return(nil, someError)

		// act
		err := f.unitRepo.Delete(ctx, unit_of_measure_id)

		// assert
		assert.ErrorIs(t, err, someError)
	})
}

func TestListUnitOfMeasures(t *testing.T) {
	ctx := context.Background()
	unit_of_measure := &models.UnitOfMeasure{Name: "UnitOfMeasure"}
	unit_of_measures := []*models.UnitOfMeasure{unit_of_measure}
	columns := []string{"unit_of_measure_id", "name"}
	args := make([]interface{}, 0)

	t.Run("list unit_of_measures", func(t *testing.T) {
		// arrange
		f := setUp(t)
		defer f.tearDown()

		queryList := "SELECT unit_of_measure_id, name FROM unit_of_measure;"
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(unit_of_measure.UnitOfMeasureId, unit_of_measure.Name).ToPgxRows()
		f.pool.EXPECT().Query(ctx, queryList, args).Return(pgxRows, nil)

		// act
		result, err := f.unitRepo.List(ctx)

		// assert
		require.NoError(t, err)
		assert.Equal(t, unit_of_measures, result)
	})
}
