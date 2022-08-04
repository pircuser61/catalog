package grpc

import (
	"context"
	"errors"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	countryPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) CountryCreate(ctx context.Context, in *pb.CountryCreateRequest) (*emptypb.Empty, error) {
	if err := i.country.Add(ctx, &models.Country{
		Name: in.GetName(),
	}); err != nil {
		if errors.Is(err, models.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) CountryUpdate(ctx context.Context, in *pb.CountryUpdateRequest) (*emptypb.Empty, error) {
	inCountry := in.GetCountry()
	if inCountry == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if err := i.country.Update(ctx, &models.Country{
		CountryId: inCountry.GetCountryId(),
		Name:      inCountry.GetName(),
	}); err != nil {
		if errors.Is(err, models.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, countryPkg.ErrCountryNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) CountryDelete(ctx context.Context, in *pb.CountryDeleteRequest) (*emptypb.Empty, error) {
	if err := i.country.Delete(ctx, uint32(in.GetCountryId())); err != nil {
		if errors.Is(err, models.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, countryPkg.ErrCountryNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) CountryList(ctx context.Context, _ *emptypb.Empty) (*pb.CountryListResponse, error) {
	countries, err := i.country.List(ctx)
	if err != nil {
		if errors.Is(err, storePkg.ErrTimeout) {
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

func (i *Implementation) CountryGet(ctx context.Context, in *pb.CountryGetRequest) (*pb.CountryGetResponse, error) {
	country, err := i.country.Get(ctx, in.GetCountryId())
	if err != nil {
		if errors.Is(err, countryPkg.ErrCountryNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
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
