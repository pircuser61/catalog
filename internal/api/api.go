package api

import (
	pb "gitlab.ozon.dev/pircuser61/catalog/api"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
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
