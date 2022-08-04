package grpc

import (
	"context"

	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	countryPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	unitOfMeasurePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type Implementation struct {
	pb.UnimplementedCatalogServer //	mustEmbedUnimplementedCatalogServer()
	good                          goodPkg.Interface
	country                       countryPkg.Interface
	unitOfMeasure                 unitOfMeasurePkg.Interface
}

func New(ctx context.Context, store storePkg.Interface) pb.CatalogServer {
	core := store.GetCore(ctx)

	return &Implementation{
		good:          core.Good,
		country:       core.Country,
		unitOfMeasure: core.UnitOfMeasure,
	}
}
