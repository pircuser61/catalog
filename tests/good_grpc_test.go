//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type goodRow struct {
	Code            uint64
	Name            string
	UnitOfMeasureId uint32
	CountryId       uint32
}

func TestGoodCreate(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		//arrange
		Db.SetUp(ctx, t)
		defer Db.TearDown(ctx)
		name := "Товар1"
		unit, country := Db.GetFirstUnitOfMeasure(ctx, t), Db.GetFirstCountry(ctx, t)
		grpcRequest := &pb.GoodCreateRequest{Name: name, UnitOfMeasure: unit.Name, Country: country.Name}
		//act
		_, err := CatalogClient.GoodCreate(ctx, grpcRequest)

		//assert
		require.NoError(t, err)
		createdRow := mustOneGood(ctx, t)
		assert.Equal(t, name, createdRow.Name)
		assert.Equal(t, unit.Id, createdRow.UnitOfMeasureId)
		assert.Equal(t, country.Id, createdRow.CountryId)
	})
}

func TestGoodGet(t *testing.T) {
	ctx := context.Background()

	t.Run("grpc good get: success", func(t *testing.T) {
		//arrange
		Db.SetUp(ctx, t)
		defer Db.TearDown(ctx)

		code, name := uint64(1), "Товар1"
		unit, country := Db.GetFirstUnitOfMeasure(ctx, t), Db.GetFirstCountry(ctx, t)
		createGood(ctx, t, code, name, unit.Id, country.Id)

		grpcRequest := &pb.GoodGetRequest{Code: code}

		//act
		result, err := CatalogClient.GoodGet(ctx, grpcRequest)

		//assert
		require.NoError(t, err)
		assert.Equal(t, code, result.Good.Code)
		assert.Equal(t, name, result.Good.Name)
		assert.Equal(t, unit.Name, result.Good.UnitOfMeasure)
		assert.Equal(t, country.Name, result.Good.Country)
	})

	t.Run("grpc good get: not found", func(t *testing.T) {
		//arrange
		Db.SetUp(ctx, t)
		defer Db.TearDown(ctx)

		grpcRequest := &pb.GoodGetRequest{Code: 1}

		//act
		_, err := CatalogClient.GoodGet(ctx, grpcRequest)

		//assert
		assert.ErrorIs(t, err, status.Error(codes.NotFound, "obj does not exist"))
	})
}

func TestGoodUpdate(t *testing.T) {
	ctx := context.Background()

	t.Run("grpc good get: success", func(t *testing.T) {
		//arrange
		Db.SetUp(ctx, t)
		defer Db.TearDown(ctx)

		code := uint64(1)
		unit, country := Db.GetFirstUnitOfMeasure(ctx, t), Db.GetFirstCountry(ctx, t)
		createGood(ctx, t, code, "Товар1", unit.Id, country.Id)
		newName := "NewName"
		grpcRequest := &pb.GoodUpdateRequest{Good: &pb.Good{
			Code:          code,
			Name:          newName,
			UnitOfMeasure: unit.Name,
			Country:       country.Name,
		}}

		//act
		_, err := CatalogClient.GoodUpdate(ctx, grpcRequest)

		//assert
		require.NoError(t, err)
		updatedRow := mustOneGood(ctx, t)
		assert.Equal(t, code, updatedRow.Code)
		assert.Equal(t, newName, updatedRow.Name)
		assert.Equal(t, unit.Id, updatedRow.UnitOfMeasureId)
		assert.Equal(t, country.Id, updatedRow.CountryId)
	})

	t.Run("grpc good update: not found", func(t *testing.T) {
		//arrange
		Db.SetUp(ctx, t)
		defer Db.TearDown(ctx)
		unit, country := Db.GetFirstUnitOfMeasure(ctx, t), Db.GetFirstCountry(ctx, t)
		grpcRequest := &pb.GoodUpdateRequest{Good: &pb.Good{
			Code:          1,
			UnitOfMeasure: unit.Name,
			Country:       country.Name,
		}}
		//act
		_, err := CatalogClient.GoodUpdate(ctx, grpcRequest)

		//assert
		assert.ErrorIs(t, err, status.Error(codes.NotFound, "obj does not exist"))
	})
}

func mustOneGood(ctx context.Context, t *testing.T) *goodRow {
	var result []*goodRow
	err := pgxscan.Select(ctx, Db.Pool, &result, "SELECT code, name, unit_of_measure_id, country_id FROM good")
	require.NoError(t, err)
	require.Equal(t, 1, len(result))
	return result[0]
}

func createGood(ctx context.Context, t *testing.T, code uint64, name string, unit_id, country_id uint32) {
	queryAdd := "INSERT INTO good (code, name, unit_of_measure_id, country_id) VALUES ($1, $2, $3, $4);"
	_, err := Db.Pool.Exec(ctx, queryAdd, code, name, unit_id, country_id)
	if err != nil {
		panic(err)
	}
}
