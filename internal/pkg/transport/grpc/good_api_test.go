package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	goodFxtr "gitlab.ozon.dev/pircuser61/catalog/tests/fixtures"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewGoodService(t *testing.T) {
	t.Run("good add success", func(t *testing.T) {
		// arrange
		good := goodFxtr.Good().Name("name").UnitOfMeasure("uom").Country("country").P()
		one := 1
		id := uint32(1)
		keys := &goodPkg.GoodKeys{&one, &id, &id}
		goodCreateMsg := pb.GoodCreateRequest{
			Name:          good.Name,
			UnitOfMeasure: good.UnitOfMeasure,
			Country:       good.Country,
		}

		f := catalogSetUp(t)
		f.goodRepo.EXPECT().Add(f.Ctx, good, keys).Return(nil)
		f.goodRepo.EXPECT().GetKeys(f.Ctx, good).Return(keys, nil)

		// act
		resp, err := f.grpcImplementation.GoodCreate(context.Background(), &goodCreateMsg)

		// assert
		require.NoError(t, err)
		assert.Equal(t, resp, &pb.GoodCreateResponse{})
	})
	t.Run("error ", func(t *testing.T) {
		t.Run("keys not found", func(t *testing.T) {
			// arrange
			good := goodFxtr.Good().Name("").UnitOfMeasure("uom").Country("country").P()
			goodCreateMsg := pb.GoodCreateRequest{
				Name:          good.Name,
				UnitOfMeasure: good.UnitOfMeasure,
				Country:       good.Country,
			}
			f := catalogSetUp(t)
			f.goodRepo.EXPECT().GetKeys(f.Ctx, good).Return(nil, models.ErrValidation)
			awaitError := status.Error(codes.InvalidArgument, "invalid data")
			// act
			_, err := f.grpcImplementation.GoodCreate(context.Background(), &goodCreateMsg)

			// assert
			assert.ErrorIs(t, err, awaitError)
		})

	})
}
