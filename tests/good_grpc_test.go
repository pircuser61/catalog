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
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
)

func TestGoodCreate(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		//arrange
		Db.SetUp(ctx, t)
		defer Db.TearDown(ctx)
		unit, country := Db.GetKeys(ctx, t)
		grpcRequest := &pb.GoodCreateRequest{Name: "Товар1", UnitOfMeasure: unit.Name, Country: country.Name}
		//act
		_, err := CatalogClient.GoodCreate(ctx, grpcRequest)

		//assert
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
		assert.Equal(t, "Товар1", createdRow.Name)
		assert.Equal(t, unit.Id, createdRow.UnitOfMeasureId)
		assert.Equal(t, country.Id, createdRow.CountryId)
	})
}
