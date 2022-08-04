package grpc

import (
	"context"
	"errors"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) GoodCreate(ctx context.Context, in *pb.GoodCreateRequest) (*emptypb.Empty, error) {
	if err := i.good.Add(ctx, &models.Good{
		Name:          in.GetName(),
		UnitOfMeasure: in.GetUnitOfMeasure(),
		Country:       in.GetCountry(),
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

func (i *Implementation) GoodUpdate(ctx context.Context, in *pb.GoodUpdateRequest) (*emptypb.Empty, error) {
	inGood := in.GetGood()
	if inGood == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if err := i.good.Update(ctx, &models.Good{
		Code:          inGood.GetCode(),
		Name:          inGood.GetName(),
		UnitOfMeasure: inGood.GetUnitOfMeasure(),
		Country:       inGood.GetCountry()}); err != nil {
		if errors.Is(err, models.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, goodPkg.ErrGoodNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) GoodDelete(ctx context.Context, in *pb.GoodDeleteRequest) (*emptypb.Empty, error) {
	if err := i.good.Delete(ctx, in.GetCode()); err != nil {
		if errors.Is(err, models.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, goodPkg.ErrGoodNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) GoodList(ctx context.Context, request *pb.GoodListRequest) (*pb.GoodListResponse, error) {
	goods, err := i.good.ListEx(ctx, request.GetLimit(), request.GetOffset())
	if err != nil {
		if errors.Is(err, storePkg.ErrTimeout) {
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

func (i *Implementation) GoodGet(ctx context.Context, in *pb.GoodGetRequest) (*pb.GoodGetResponse, error) {
	good, err := i.good.Get(ctx, in.GetCode())
	if err != nil {
		if errors.Is(err, goodPkg.ErrGoodNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, storePkg.ErrTimeout) {
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
