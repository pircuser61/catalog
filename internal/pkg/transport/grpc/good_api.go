package grpc

import (
	"context"
	"errors"
	"io"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/models"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GoodCreate(stream pb.Catalog_GoodCreateServer) error {

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		response := pb.GoodCreateResponse{}
		ctx := context.Background()
		if err := i.good.Add(ctx, &models.Good{
			Name:          in.GetName(),
			UnitOfMeasure: in.GetUnitOfMeasure(),
			Country:       in.GetCountry(),
		}); err != nil {
			if errors.Is(err, models.ErrValidation) {
				err = status.Error(codes.InvalidArgument, err.Error())
			} else if errors.Is(err, storePkg.ErrTimeout) {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			} else {
				err = status.Error(codes.Internal, err.Error())
			}
			errMsg := err.Error()
			response.Error = &errMsg
		}

		if err := stream.Send(&response); err != nil {
			return err
		}
	}
}

func (i *Implementation) GoodUpdate(stream pb.Catalog_GoodUpdateServer) error {

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		response := pb.GoodUpdateResponse{}
		ctx := context.Background()
		inGood := in.Good
		if err := i.good.Update(ctx, &models.Good{
			Code:          inGood.GetCode(),
			Name:          inGood.GetName(),
			UnitOfMeasure: inGood.GetUnitOfMeasure(),
			Country:       inGood.GetCountry()}); err != nil {
			if errors.Is(err, models.ErrValidation) {
				err = status.Error(codes.InvalidArgument, err.Error())
			} else if errors.Is(err, storePkg.ErrTimeout) {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			} else if errors.Is(err, goodPkg.ErrGoodNotFound) {
				err = status.Error(codes.NotFound, err.Error())
			} else {
				err = status.Error(codes.Internal, err.Error())
			}
			errMsg := err.Error()
			response.Error = &errMsg
		}

		if err := stream.Send(&response); err != nil {
			return err
		}
	}
}

func (i *Implementation) GoodDelete(stream pb.Catalog_GoodDeleteServer) error {

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		result := pb.GoodDeleteResponse{}
		ctx := context.Background()
		if err := i.good.Delete(ctx, in.GetCode()); err != nil {
			if errors.Is(err, goodPkg.ErrGoodNotFound) {
				err = status.Error(codes.NotFound, err.Error())
			}
			if errors.Is(err, storePkg.ErrTimeout) {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			}
			errMsg := err.Error()
			result.Error = &errMsg
		}
		if err := stream.Send(&result); err != nil {
			return err
		}
	}
}

func (i *Implementation) GoodDelete2(ctx context.Context, in *pb.GoodDeleteRequest) (*pb.GoodDeleteResponse, error) {
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
	return &pb.GoodDeleteResponse{}, nil
}

func (i *Implementation) GoodList(stream pb.Catalog_GoodListServer) error {

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		result := &pb.GoodListResponse{}
		ctx := context.Background()
		goods, err := i.good.List(ctx, in.GetLimit(), in.GetOffset())
		if err != nil {
			if errors.Is(err, storePkg.ErrTimeout) {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			} else {
				err = status.Error(codes.Internal, err.Error())
			}
			errMsg := err.Error()
			result.Error = &errMsg
		} else {
			listGoods := make([]*pb.GoodListResponse_Good, 0, len(goods))
			for _, good := range goods {
				listGoods = append(listGoods, &pb.GoodListResponse_Good{
					Code: good.Code,
					Name: good.Name,
				})
			}
			result.Goods = listGoods
		}

		if err := stream.Send(result); err != nil {
			return err
		}
	}
}

func (i *Implementation) GoodGet(stream pb.Catalog_GoodGetServer) error {

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		ctx := context.Background()
		good, err := i.good.Get(ctx, in.GetCode())
		result := pb.GoodGetResponse{
			Good: &pb.Good{},
		}

		if err != nil {

			if errors.Is(err, goodPkg.ErrGoodNotFound) {
				err = status.Error(codes.NotFound, err.Error())
			}
			if errors.Is(err, storePkg.ErrTimeout) {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			}
			errMsg := err.Error()
			result.Error = &errMsg
		} else {
			result.Good.Code = good.Code
			result.Good.Name = good.Name
			result.Good.UnitOfMeasure = good.UnitOfMeasure
			result.Good.Country = good.Country
		}
		if err := stream.Send(&result); err != nil {
			return err
		}
	}
}
