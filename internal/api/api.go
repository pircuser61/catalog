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

type implementation struct {
	pb.UnimplementedCatalogServer //	mustEmbedUnimplementedCatalogServer()
	good                          goodPkg.Interface
}

func New(good goodPkg.Interface) pb.CatalogServer {
	return &implementation{
		good: good,
	}
}

func (i *implementation) GoodCreate(ctx context.Context, in *pb.GoodCreateRequest) (*emptypb.Empty, error) {
	if err := i.good.GoodCreate(ctx, models.Good{
		Name:          in.GetName(),
		UnitOfMeasure: in.GetUnitOfMeasure(),
		Country:       in.GetCountry(),
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

func (i *implementation) GoodUpdate(ctx context.Context, in *pb.GoodUpdateRequest) (*emptypb.Empty, error) {
	inGood := in.GetGood()
	if inGood == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if err := i.good.GoodUpdate(ctx, models.Good{
		Code:          inGood.GetCode(),
		Name:          inGood.GetName(),
		UnitOfMeasure: inGood.GetUnitOfMeasure(),
		Country:       inGood.GetCountry()}); err != nil {
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

func (i *implementation) GoodDelete(ctx context.Context, in *pb.GoodDeleteRequest) (*emptypb.Empty, error) {
	if err := i.good.GoodDelete(ctx, in.GetCode()); err != nil {
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

func (i *implementation) GoodList(ctx context.Context, _ *emptypb.Empty) (*pb.GoodListResponse, error) {
	goods, err := i.good.GoodList(ctx)
	if err != nil {
		if errors.Is(err, cache.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := make([]*pb.GoodListResponse_Good, 0, len(goods))
	for _, good := range goods {
		result = append(result, &pb.GoodListResponse_Good{
			Code: good.Code,
			Name: good.Name,
		})
	}
	return &pb.GoodListResponse{
		Goods: result,
	}, nil
}

func (i *implementation) GoodGet(ctx context.Context, in *pb.GoodGetRequest) (*pb.GoodGetResponse, error) {
	good, err := i.good.GoodGet(ctx, in.GetCode())
	if err != nil {
		if errors.Is(err, goodPkg.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, cache.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GoodGetResponse{
		Good: &pb.Good{
			Code:          good.Code,
			Name:          good.Name,
			UnitOfMeasure: good.UnitOfMeasure,
			Country:       good.Country},
	}, nil
}
