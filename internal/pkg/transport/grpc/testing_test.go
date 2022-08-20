package grpc

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	mockGoodRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/mocks"
	goodUseCase "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/usecase"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type serviceFixture struct {
	Ctx                context.Context
	grpcImplementation pb.CatalogServer
}

type goodFixture struct {
	Ctx                context.Context
	goodRepo           *mockGoodRepo.MockRepository
	grpcImplementation pb.CatalogServer
}

/*
func serviceSetUp(t *testing.T) serviceFixture {
	t.Parallel()
	f := serviceFixture{Ctx: context.Background()}
	f.grpcImplementation = New(f.Ctx, f.core)
	return f
}
*/

func catalogSetUp(t *testing.T) goodFixture {
	t.Parallel()

	f := goodFixture{}
	f.Ctx = context.Background()
	f.goodRepo = mockGoodRepo.NewMockRepository(gomock.NewController(t))

	store := &storePkg.Core{
		Good:          goodUseCase.New(f.goodRepo),
		Country:       nil,
		UnitOfMeasure: nil,
	}

	f.grpcImplementation = New(f.Ctx, store)
	return f
}
