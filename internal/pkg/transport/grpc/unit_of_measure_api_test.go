package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func getSomeUnitOfMeasure(id uint32) *models.UnitOfMeasure {
	return &models.UnitOfMeasure{UnitOfMeasureId: id, Name: "county_name"}
}

func TestNewUnitOfMeasureService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func(unit *models.UnitOfMeasure) *pb.UnitOfMeasureCreateRequest {
		return &pb.UnitOfMeasureCreateRequest{
			Name: unit.Name,
		}
	}

	t.Run("unit add success", func(t *testing.T) {
		// arrange
		unit := getSomeUnitOfMeasure(0)
		reqMsg := getReqMsg(unit)

		f := catalogSetUp(ctx, t)
		f.unitRepo.EXPECT().Add(ctx, unit).Return(nil)

		// act
		resp, err := f.grpcImplementation.UnitOfMeasureCreate(ctx, reqMsg)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &emptypb.Empty{}, resp)
	})

	t.Run("create error ", func(t *testing.T) {
		cases := []struct {
			name   string
			errIn  error
			errOut error
		}{
			{"timeout", storePkg.ErrTimeout, status.Error(codes.DeadlineExceeded, "Timeout")},
			{"internal error", errors.New("error text"), status.Error(codes.Internal, "error text")},
		}

		for _, tc := range cases {
			itc := tc
			t.Run(itc.name, func(t *testing.T) {
				// arrange
				unit := getSomeUnitOfMeasure(0)
				reqMsg := getReqMsg(unit)
				f := catalogSetUp(ctx, t)
				f.unitRepo.EXPECT().Add(ctx, unit).Return(itc.errIn)

				// act
				_, err := f.grpcImplementation.UnitOfMeasureCreate(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}

func TestGetUnitOfMeasureService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func() *pb.UnitOfMeasureGetRequest {
		return &pb.UnitOfMeasureGetRequest{
			UnitOfMeasureId: 1,
		}
	}

	t.Run("unit get success", func(t *testing.T) {
		// arrange
		reqMsg := getReqMsg()
		unit := getSomeUnitOfMeasure(1)
		f := catalogSetUp(ctx, t)
		f.unitRepo.EXPECT().Get(ctx, reqMsg.UnitOfMeasureId).Return(unit, nil)

		// act
		resp, err := f.grpcImplementation.UnitOfMeasureGet(ctx, reqMsg)

		// assert
		respMsg := &pb.UnitOfMeasureGetResponse{
			Unit: &pb.UnitOfMeasure{
				UnitOfMeasureId: unit.UnitOfMeasureId,
				Name:            unit.Name},
		}
		require.NoError(t, err)
		assert.Equal(t, respMsg, resp)
	})

	t.Run("error ", func(t *testing.T) {
		cases := []struct {
			name   string
			errIn  error
			errOut error
		}{{"not found", storePkg.ErrNotExists, status.Error(codes.NotFound, "obj does not exist")},
			{"timeout", storePkg.ErrTimeout, status.Error(codes.DeadlineExceeded, "Timeout")},
			{"internal error", errors.New("error text"), status.Error(codes.Internal, "error text")}}

		for _, tc := range cases {
			itc := tc
			t.Run(itc.name, func(t *testing.T) {
				// arrange
				reqMsg := getReqMsg()
				f := catalogSetUp(ctx, t)
				f.unitRepo.EXPECT().Get(ctx, reqMsg.UnitOfMeasureId).Return(nil, itc.errIn)

				// act
				_, err := f.grpcImplementation.UnitOfMeasureGet(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}

func TestUpdateUnitOfMeasureService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func(unit *models.UnitOfMeasure) *pb.UnitOfMeasureUpdateRequest {
		if unit == nil {
			return &pb.UnitOfMeasureUpdateRequest{}
		}
		return &pb.UnitOfMeasureUpdateRequest{
			Unit: &pb.UnitOfMeasure{
				UnitOfMeasureId: unit.UnitOfMeasureId,
				Name:            unit.Name},
		}
	}

	t.Run("unit update success", func(t *testing.T) {
		// arrange
		unit := getSomeUnitOfMeasure(1)
		reqMsg := getReqMsg(unit)

		f := catalogSetUp(ctx, t)
		f.unitRepo.EXPECT().Update(ctx, unit).Return(nil)
		// act
		resp, err := f.grpcImplementation.UnitOfMeasureUpdate(ctx, reqMsg)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &emptypb.Empty{}, resp)
	})

	t.Run("unit update bad request", func(t *testing.T) {
		// arrange
		reqMsg := getReqMsg(nil)
		f := catalogSetUp(ctx, t)

		// act
		_, err := f.grpcImplementation.UnitOfMeasureUpdate(ctx, reqMsg)

		// assert
		assert.ErrorIs(t, status.Error(codes.InvalidArgument, "empty request"), err)
	})

	t.Run("update error ", func(t *testing.T) {
		cases := []struct {
			name   string
			errIn  error
			errOut error
		}{{"not found", storePkg.ErrNotExists, status.Error(codes.NotFound, "obj does not exist")},
			{"timeout", storePkg.ErrTimeout, status.Error(codes.DeadlineExceeded, "Timeout")},
			{"internal error", errors.New("error text"), status.Error(codes.Internal, "error text")}}

		for _, tc := range cases {
			itc := tc
			t.Run(itc.name, func(t *testing.T) {
				// arrange
				unit := getSomeUnitOfMeasure(0)
				reqMsg := getReqMsg(unit)
				f := catalogSetUp(ctx, t)
				f.unitRepo.EXPECT().Update(ctx, unit).Return(itc.errIn)

				// act
				_, err := f.grpcImplementation.UnitOfMeasureUpdate(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}

func TestUnitDeleteService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func() *pb.UnitOfMeasureDeleteRequest {
		return &pb.UnitOfMeasureDeleteRequest{
			UnitOfMeasureId: 1,
		}
	}

	t.Run("unit delete success", func(t *testing.T) {
		// arrange
		reqMsg := getReqMsg()
		f := catalogSetUp(ctx, t)
		f.unitRepo.EXPECT().Delete(ctx, reqMsg.UnitOfMeasureId).Return(nil)

		// act
		resp, err := f.grpcImplementation.UnitOfMeasureDelete(ctx, reqMsg)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &emptypb.Empty{}, resp)
	})

	t.Run("unit delete error ", func(t *testing.T) {
		cases := []struct {
			name   string
			errIn  error
			errOut error
		}{{"not found", storePkg.ErrNotExists, status.Error(codes.NotFound, "obj does not exist")},
			{"timeout", storePkg.ErrTimeout, status.Error(codes.DeadlineExceeded, "Timeout")},
			{"internal error", errors.New("error text"), status.Error(codes.Internal, "error text")}}

		for _, tc := range cases {
			itc := tc
			t.Run(itc.name, func(t *testing.T) {
				// arrange
				reqMsg := getReqMsg()
				f := catalogSetUp(ctx, t)
				f.unitRepo.EXPECT().Delete(ctx, reqMsg.UnitOfMeasureId).Return(itc.errIn)

				// act
				_, err := f.grpcImplementation.UnitOfMeasureDelete(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}
