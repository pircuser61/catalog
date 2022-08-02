package api

import (
	"context"
	"errors"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *implementation) CountryCreate(ctx context.Context, in *pb.CountryCreateRequest) (*emptypb.Empty, error) {
	if err := i.good.CountryCreate(ctx, models.Country{
		Name: in.GetName(),
	}); err != nil {
		if errors.Is(err, models.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, cache.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *implementation) CountryUpdate(ctx context.Context, in *pb.CountryUpdateRequest) (*emptypb.Empty, error) {
	inCountry := in.GetCountry()
	if inCountry == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if err := i.good.CountryUpdate(ctx, models.Country{
		CountryId: inCountry.GetCountryId(),
		Name:      inCountry.GetName(),
	}); err != nil {
		if errors.Is(err, models.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, goodPkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, cache.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *implementation) CountryDelete(ctx context.Context, in *pb.CountryDeleteRequest) (*emptypb.Empty, error) {
	if err := i.good.CountryDelete(ctx, uint32(in.GetCountryId())); err != nil {
		if errors.Is(err, models.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, goodPkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, cache.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *implementation) CountryList(ctx context.Context, _ *emptypb.Empty) (*pb.CountryListResponse, error) {
	countries, err := i.good.CountryList(ctx)
	if err != nil {
		if errors.Is(err, cache.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := make([]*pb.Country, 0, len(countries))
	for _, country := range countries {
		result = append(result, &pb.Country{
			CountryId: country.CountryId,
			Name:      country.Name,
		})
	}
	return &pb.CountryListResponse{
		Countries: result,
	}, nil
}

func (i *implementation) CountryGet(ctx context.Context, in *pb.CountryGetRequest) (*pb.CountryGetResponse, error) {
	country, err := i.good.CountryGet(ctx, in.GetCountryId())
	if err != nil {
		if errors.Is(err, goodPkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, cache.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CountryGetResponse{
		Country: &pb.Country{
			CountryId: country.CountryId,
			Name:      country.Name},
	}, nil
}
