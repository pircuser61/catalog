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

func getSomeCountry(id uint32) *models.Country {
	return &models.Country{CountryId: id, Name: "county_name"}
}

func TestNewCountryService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func(country *models.Country) *pb.CountryCreateRequest {
		return &pb.CountryCreateRequest{
			Name: country.Name,
		}
	}

	t.Run("country add success", func(t *testing.T) {
		// arrange
		country := getSomeCountry(0)
		reqMsg := getReqMsg(country)

		f := catalogSetUp(ctx, t)
		f.countryRepo.EXPECT().Add(ctx, country).Return(nil)

		// act
		resp, err := f.grpcImplementation.CountryCreate(ctx, reqMsg)

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
				country := getSomeCountry(0)
				reqMsg := getReqMsg(country)
				f := catalogSetUp(ctx, t)
				f.countryRepo.EXPECT().Add(ctx, country).Return(itc.errIn)

				// act
				_, err := f.grpcImplementation.CountryCreate(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}

func TestGetCountryService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func() *pb.CountryGetRequest {
		return &pb.CountryGetRequest{
			CountryId: 1,
		}
	}

	t.Run("country get success", func(t *testing.T) {
		// arrange
		reqMsg := getReqMsg()
		country := getSomeCountry(1)
		f := catalogSetUp(ctx, t)
		f.countryRepo.EXPECT().Get(ctx, reqMsg.CountryId).Return(country, nil)

		// act
		resp, err := f.grpcImplementation.CountryGet(ctx, reqMsg)

		// assert
		respMsg := &pb.CountryGetResponse{
			Country: &pb.Country{
				CountryId: country.CountryId,
				Name:      country.Name},
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
				f.countryRepo.EXPECT().Get(ctx, reqMsg.CountryId).Return(nil, itc.errIn)

				// act
				_, err := f.grpcImplementation.CountryGet(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}

func TestUpdateCountryService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func(country *models.Country) *pb.CountryUpdateRequest {
		if country == nil {
			return &pb.CountryUpdateRequest{}
		}
		return &pb.CountryUpdateRequest{
			Country: &pb.Country{
				CountryId: country.CountryId,
				Name:      country.Name},
		}
	}

	t.Run("country update success", func(t *testing.T) {
		// arrange
		country := getSomeCountry(1)
		reqMsg := getReqMsg(country)

		f := catalogSetUp(ctx, t)
		f.countryRepo.EXPECT().Update(ctx, country).Return(nil)
		// act
		resp, err := f.grpcImplementation.CountryUpdate(ctx, reqMsg)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &emptypb.Empty{}, resp)
	})

	t.Run("country update bad request", func(t *testing.T) {
		// arrange
		reqMsg := getReqMsg(nil)
		f := catalogSetUp(ctx, t)

		// act
		_, err := f.grpcImplementation.CountryUpdate(ctx, reqMsg)

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
				country := getSomeCountry(0)
				reqMsg := getReqMsg(country)
				f := catalogSetUp(ctx, t)
				f.countryRepo.EXPECT().Update(ctx, country).Return(itc.errIn)

				// act
				_, err := f.grpcImplementation.CountryUpdate(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}

func TestCountryDeleteService(t *testing.T) {
	ctx := context.Background()
	getReqMsg := func() *pb.CountryDeleteRequest {
		return &pb.CountryDeleteRequest{
			CountryId: 1,
		}
	}

	t.Run("country delete success", func(t *testing.T) {
		// arrange
		reqMsg := getReqMsg()
		f := catalogSetUp(ctx, t)
		f.countryRepo.EXPECT().Delete(ctx, reqMsg.CountryId).Return(nil)

		// act
		resp, err := f.grpcImplementation.CountryDelete(ctx, reqMsg)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &emptypb.Empty{}, resp)
	})

	t.Run("country delete error ", func(t *testing.T) {
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
				f.countryRepo.EXPECT().Delete(ctx, reqMsg.CountryId).Return(itc.errIn)

				// act
				_, err := f.grpcImplementation.CountryDelete(ctx, reqMsg)

				// assert
				assert.ErrorIs(t, err, itc.errOut)
			})
		}
	})
}
