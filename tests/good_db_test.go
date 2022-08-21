//go:build integration
// +build integration

package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	goodRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/repository/postgre"
	goodFxtr "gitlab.ozon.dev/pircuser61/catalog/tests/fixtures"
)

func TestCreateGood(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		//arrange
		Db.SetUp(ctx, t)
		defer Db.TearDown(ctx)

		goodKeys := GetKeys(ctx, t)
		goodRepo := goodRepo.New(Db.Pool, Timeout)
		good := goodFxtr.Good().Name("name").P()

		// act
		err := goodRepo.Add(ctx, good, goodKeys)

		// assert
		require.NoError(t, err)
		type row struct {
			Name            string
			UnitOfMeasureId uint32
			CountryId       uint32
		}

		var result []*row
		err = pgxscan.Select(ctx, Db.Pool, &result, "SELECT name, unit_of_measure_id, country_id FROM good")
		if err != nil {
			fmt.Println(err.Error())
		}
		require.NoError(t, err)
		require.Equal(t, 1, len(result))
		createdRow := result[0]
		assert.Equal(t, good.Name, createdRow.Name)
		assert.Equal(t, *goodKeys.CountryId, createdRow.CountryId)
		assert.Equal(t, *goodKeys.UnitOfMeasureId, createdRow.UnitOfMeasureId)
	})
}

func GetKeys(ctx context.Context, t *testing.T) *goodPkg.GoodKeys {
	var unit_of_measure_id, country_id uint32
	if err := pgxscan.Get(ctx, Db.Pool, &country_id, "SELECT country_id FROM country LIMIT 1"); err != nil {
		panic("Cant get country id: " + err.Error())
	}
	if err := pgxscan.Get(ctx, Db.Pool, &unit_of_measure_id, "SELECT unit_of_measure_id FROM unit_of_measure LIMIT 1"); err != nil {
		panic("Cant get uom id: " + err.Error())
	}
	return &goodPkg.GoodKeys{UnitOfMeasureId: &unit_of_measure_id, CountryId: &country_id}
}
