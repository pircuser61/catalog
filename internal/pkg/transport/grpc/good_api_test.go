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
	goodFxtr "gitlab.ozon.dev/pircuser61/catalog/tests/fixtures"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getSomeGood(code uint64) *models.Good {
	return goodFxtr.Good().Code(code).Name("name").UnitOfMeasure("uom").Country("country").P()
}

func TestNewGoodService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func(good *models.Good) *pb.GoodCreateRequest {
		return &pb.GoodCreateRequest{
			Name:          good.Name,
			UnitOfMeasure: good.UnitOfMeasure,
			Country:       good.Country,
		}
	}

	t.Run("good add success", func(t *testing.T) {
		// arrange
		good := getSomeGood(0)
		reqMsg := getReqMsg(good)

		f := catalogSetUp(ctx, t)
		f.goodRepo.EXPECT().Add(ctx, good).Return(nil)

		// act
		resp, err := f.grpcImplementation.GoodCreate(ctx, reqMsg)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &pb.GoodCreateResponse{}, resp)
	})
	t.Run("create error ", func(t *testing.T) {

		cases := []struct {
			name   string
			errIn  error
			errOut error
		}{{"keys not found", models.ErrValidation, status.Error(codes.InvalidArgument, "invalid data")},
			{"timeout", storePkg.ErrTimeout, status.Error(codes.DeadlineExceeded, "Timeout")},
			{"internal error", errors.New("error text"), status.Error(codes.Internal, "error text")}}

		for _, tc := range cases {
			itc := tc
			t.Run(itc.name, func(t *testing.T) {
				// arrange
				good := getSomeGood(0)
				reqMsg := getReqMsg(good)
				f := catalogSetUp(ctx, t)
				f.goodRepo.EXPECT().Add(ctx, good).Return(itc.errIn)

				// act
				_, err := f.grpcImplementation.GoodCreate(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}

func TestGetGoodService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func() *pb.GoodGetRequest {
		return &pb.GoodGetRequest{
			Code: 1,
		}
	}
	t.Run("good get success", func(t *testing.T) {
		// arrange
		reqMsg := getReqMsg()
		good := getSomeGood(1)
		f := catalogSetUp(ctx, t)
		f.goodRepo.EXPECT().Get(ctx, reqMsg.Code).Return(good, nil)

		// act
		resp, err := f.grpcImplementation.GoodGet(ctx, reqMsg)

		// assert
		respMsg := &pb.GoodGetResponse{
			Good: &pb.Good{
				Code:          good.Code,
				Name:          good.Name,
				UnitOfMeasure: good.UnitOfMeasure,
				Country:       good.Country},
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
				f.goodRepo.EXPECT().Get(ctx, reqMsg.Code).Return(nil, itc.errIn)

				// act
				_, err := f.grpcImplementation.GoodGet(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}

func TestUpdateGoodService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func(good *models.Good) *pb.GoodUpdateRequest {
		if good == nil {
			return &pb.GoodUpdateRequest{}
		}
		return &pb.GoodUpdateRequest{
			Good: &pb.Good{
				Code:          good.Code,
				Name:          good.Name,
				UnitOfMeasure: good.UnitOfMeasure,
				Country:       good.Country},
		}
	}

	t.Run("good update success", func(t *testing.T) {
		// arrange
		good := getSomeGood(1)
		reqMsg := getReqMsg(good)

		f := catalogSetUp(ctx, t)
		f.goodRepo.EXPECT().Update(ctx, good).Return(nil)
		// act
		resp, err := f.grpcImplementation.GoodUpdate(ctx, reqMsg)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &pb.GoodUpdateResponse{}, resp)
	})

	t.Run("good update bad request", func(t *testing.T) {
		// arrange
		reqMsg := getReqMsg(nil)
		f := catalogSetUp(ctx, t)

		// act
		_, err := f.grpcImplementation.GoodUpdate(ctx, reqMsg)

		// assert
		assert.ErrorIs(t, status.Error(codes.InvalidArgument, "empty request"), err)
	})

	t.Run("update error ", func(t *testing.T) {
		cases := []struct {
			name   string
			errIn  error
			errOut error
		}{{"not found", storePkg.ErrNotExists, status.Error(codes.NotFound, "obj does not exist")},
			{"keys not found", models.ErrValidation, status.Error(codes.InvalidArgument, "invalid data")},
			{"timeout", storePkg.ErrTimeout, status.Error(codes.DeadlineExceeded, "Timeout")},
			{"internal error", errors.New("error text"), status.Error(codes.Internal, "error text")}}

		for _, tc := range cases {
			itc := tc
			t.Run(itc.name, func(t *testing.T) {
				// arrange
				good := getSomeGood(0)
				reqMsg := getReqMsg(good)
				f := catalogSetUp(ctx, t)
				f.goodRepo.EXPECT().Update(ctx, good).Return(itc.errIn)

				// act
				_, err := f.grpcImplementation.GoodUpdate(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}

func TestGoodDeleteService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func() *pb.GoodDeleteRequest {
		return &pb.GoodDeleteRequest{
			Code: 1,
		}
	}

	t.Run("good delete success", func(t *testing.T) {
		// arrange
		reqMsg := getReqMsg()
		f := catalogSetUp(ctx, t)
		f.goodRepo.EXPECT().Delete(ctx, reqMsg.Code).Return(nil)

		// act
		resp, err := f.grpcImplementation.GoodDelete(ctx, reqMsg)

		// assert
		respMsg := &pb.GoodDeleteResponse{}
		require.NoError(t, err)
		assert.Equal(t, respMsg, resp)
	})

	t.Run("good delete error ", func(t *testing.T) {
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
				f.goodRepo.EXPECT().Delete(ctx, reqMsg.Code).Return(itc.errIn)

				// act
				_, err := f.grpcImplementation.GoodDelete(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}

func TestGetGoodList(t *testing.T) {
	ctx := context.Background()
	t.Run("good list sucsess", func(t *testing.T) {
		//arrange
		f := catalogSetUp(ctx, t)
		good := getSomeGood(1)
		listGoods := []*models.Good{good}
		var x uint64
		f.goodRepo.EXPECT().List(ctx, x, x).Return(listGoods, nil)

		// act
		respMsg := getGoodList(f.goodRepo, x, x)

		// assert
		listItem := &pb.GoodListResponse_Good{Name: good.Name, Code: good.Code}
		awaitResp := &pb.GoodListResponse{Goods: []*pb.GoodListResponse_Good{
			listItem}}
		assert.Equal(t, awaitResp, respMsg)
	})
	t.Run("good list error ", func(t *testing.T) {
		cases := []struct {
			name   string
			errIn  error
			errOut error
		}{
			{"timeout", storePkg.ErrTimeout, status.Error(codes.DeadlineExceeded, "Timeout")},
			{"internal error", errors.New("error text"), status.Error(codes.Internal, "error text")}}

		for _, tc := range cases {
			itc := tc
			t.Run(itc.name, func(t *testing.T) {
				// arrange
				f := catalogSetUp(ctx, t)

				var x uint64
				f.goodRepo.EXPECT().List(ctx, x, x).Return(nil, itc.errIn)

				// act
				respMsg := getGoodList(f.goodRepo, x, x)

				// assert
				errMsg := itc.errOut.Error()
				assert.Equal(t, &pb.GoodListResponse{Error: &errMsg}, respMsg)
			})
		}
	})
}
