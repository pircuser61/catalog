package grpc

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	mockCountryRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country/mocks"
	mockGoodRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/mocks"
	mockUnitRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure/mocks"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type catalogFixture struct {
	goodRepo           *mockGoodRepo.MockRepository
	countryRepo        *mockCountryRepo.MockRepository
	unitRepo           *mockUnitRepo.MockRepository
	grpcImplementation pb.CatalogServer
}

func catalogSetUp(ctx context.Context, t *testing.T) catalogFixture {
	t.Parallel()

	f := catalogFixture{}
	ctrl := gomock.NewController(t)

	f.goodRepo = mockGoodRepo.NewMockRepository(ctrl)
	f.countryRepo = mockCountryRepo.NewMockRepository(ctrl)
	f.unitRepo = mockUnitRepo.NewMockRepository(ctrl)

	store := &storePkg.Core{
		Good:          f.goodRepo,
		Country:       f.countryRepo,
		UnitOfMeasure: f.unitRepo,
	}

	f.grpcImplementation = New(ctx, store)
	return f
}
