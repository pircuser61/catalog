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

func (i *implementation) GoodCreate(_ context.Context, in *pb.GoodCreateRequest) (*pb.GoodCreateResponse, error) {
	if err := i.good.Create(models.Good{
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
	return &pb.GoodCreateResponse{}, nil
}

func (i *implementation) GoodUpdate(_ context.Context, in *pb.GoodUpdateRequest) (*pb.GoodUpdateResponse, error) {
	if err := i.good.Update(models.Good{
		Code:          in.GetCode(),
		Name:          in.GetName(),
		UnitOfMeasure: in.GetUnitOfMeasure(),
		Country:       in.GetCountry()}); err != nil {
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
	return &pb.GoodUpdateResponse{}, nil
}

func (i *implementation) GoodDelete(_ context.Context, in *pb.GoodDeleteRequest) (*pb.GoodDeleteResponse, error) {
	if err := i.good.Delete(in.GetCode()); err != nil {
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
	return &pb.GoodDeleteResponse{}, nil
}

func (i *implementation) GoodList(context.Context, *pb.GoodListRequest) (*pb.GoodListResponse, error) {
	goods, err := i.good.List()
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

func (i *implementation) GoodGet(_ context.Context, in *pb.GoodGetRequest) (*pb.GoodGetResponse, error) {
	good, err := i.good.Get(in.GetCode())
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
		Code:          good.Code,
		Name:          good.Name,
		UnitOfMeasure: good.UnitOfMeasure,
		Country:       good.Country,
	}, nil
}
