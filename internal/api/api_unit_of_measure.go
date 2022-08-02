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

func (i *implementation) UnitOfMeasureCreate(ctx context.Context, in *pb.UnitOfMeasureCreateRequest) (*emptypb.Empty, error) {
	if err := i.good.UnitOfMeasureCreate(ctx, models.UnitOfMeasure{
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

func (i *implementation) UnitOfMeasureUpdate(ctx context.Context, in *pb.UnitOfMeasureUpdateRequest) (*emptypb.Empty, error) {
	inUnitOfMeasure := in.GetUnit()
	if inUnitOfMeasure == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if err := i.good.UnitOfMeasureUpdate(ctx, models.UnitOfMeasure{
		UnitOfMeasureId: inUnitOfMeasure.GetUnitOfMeasureId(),
		Name:            inUnitOfMeasure.GetName(),
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

func (i *implementation) UnitOfMeasureDelete(ctx context.Context, in *pb.UnitOfMeasureDeleteRequest) (*emptypb.Empty, error) {
	if err := i.good.UnitOfMeasureDelete(ctx, uint32(in.GetUnitOfMeasureId())); err != nil {
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

func (i *implementation) UnitOfMeasureList(ctx context.Context, _ *emptypb.Empty) (*pb.UnitOfMeasureListResponse, error) {
	countries, err := i.good.UnitOfMeasureList(ctx)
	if err != nil {
		if errors.Is(err, cache.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := make([]*pb.UnitOfMeasure, 0, len(countries))
	for _, uom := range countries {
		result = append(result, &pb.UnitOfMeasure{
			UnitOfMeasureId: uom.UnitOfMeasureId,
			Name:            uom.Name,
		})
	}
	return &pb.UnitOfMeasureListResponse{
		Units: result,
	}, nil
}

func (i *implementation) UnitOfMeasureGet(ctx context.Context, in *pb.UnitOfMeasureGetRequest) (*pb.UnitOfMeasureGetResponse, error) {
	country, err := i.good.UnitOfMeasureGet(ctx, in.GetUnitOfMeasureId())
	if err != nil {
		if errors.Is(err, goodPkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, cache.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UnitOfMeasureGetResponse{
		Unit: &pb.UnitOfMeasure{
			UnitOfMeasureId: country.UnitOfMeasureId,
			Name:            country.Name},
	}, nil
}
