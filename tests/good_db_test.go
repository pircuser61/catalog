//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goodRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/repository/postgre"
	goodFxtr "gitlab.ozon.dev/pircuser61/catalog/tests/fixtures"
)

func TestCreateGood(t *testing.T) {
	ctx := context.Background()
	timeout := time.Duration(time.Millisecond * 1000)
	t.Run("good create success", func(t *testing.T) {
		//arrange
		Db.SetUp(ctx, t)
		defer Db.TearDown(ctx)

		country := Db.GetFirstCountry(ctx, t)
		unit := Db.GetFirstUnitOfMeasure(ctx, t)
		good := goodFxtr.Good().Name("name").UnitOfMeasure(unit.Name).Country(country.Name).P()

		goodCore := goodRepo.New(Db.Pool, timeout)

		// act
		err := goodCore.Add(ctx, good)

		// assert
		require.NoError(t, err)
		type row struct {
			Name            string
			UnitOfMeasureId uint32
			CountryId       uint32
		}

		var result []*row
		err = pgxscan.Select(ctx, Db.Pool, &result, "SELECT name, unit_of_measure_id, country_id FROM good")
		require.NoError(t, err)
		require.Equal(t, 1, len(result))
		createdRow := result[0]
		assert.Equal(t, good.Name, createdRow.Name)
		assert.Equal(t, country.Id, createdRow.CountryId)
		assert.Equal(t, unit.Id, createdRow.UnitOfMeasureId)
	})
}
