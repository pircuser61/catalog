package grpc

import (
	"context"
	"errors"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) UnitOfMeasureCreate(ctx context.Context, in *pb.UnitOfMeasureCreateRequest) (*emptypb.Empty, error) {
	if err := i.unitOfMeasure.Add(ctx, &models.UnitOfMeasure{
		Name: in.GetName(),
	}); err != nil {
		if errors.Is(err, storePkg.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) UnitOfMeasureUpdate(ctx context.Context, in *pb.UnitOfMeasureUpdateRequest) (*emptypb.Empty, error) {
	inUnitOfMeasure := in.GetUnit()
	if inUnitOfMeasure == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if err := i.unitOfMeasure.Update(ctx, &models.UnitOfMeasure{
		UnitOfMeasureId: inUnitOfMeasure.GetUnitOfMeasureId(),
		Name:            inUnitOfMeasure.GetName(),
	}); err != nil {
		if errors.Is(err, storePkg.ErrNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) UnitOfMeasureDelete(ctx context.Context, in *pb.UnitOfMeasureDeleteRequest) (*emptypb.Empty, error) {
	if err := i.unitOfMeasure.Delete(ctx, uint32(in.GetUnitOfMeasureId())); err != nil {
		if errors.Is(err, storePkg.ErrNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) UnitOfMeasureList(ctx context.Context, _ *emptypb.Empty) (*pb.UnitOfMeasureListResponse, error) {
	countries, err := i.unitOfMeasure.List(ctx)
	if err != nil {
		if errors.Is(err, storePkg.ErrTimeout) {
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

func (i *Implementation) UnitOfMeasureGet(ctx context.Context, in *pb.UnitOfMeasureGetRequest) (*pb.UnitOfMeasureGetResponse, error) {
	country, err := i.unitOfMeasure.Get(ctx, in.GetUnitOfMeasureId())
	if err != nil {
		if errors.Is(err, storePkg.ErrNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
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
