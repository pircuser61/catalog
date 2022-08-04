package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/pircuser61/catalog/internal/config"

	//countryPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country"
	countryRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country/repository/postgre"
	countryUseCase "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/country/usecase"

	//	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	goodRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/repository/postgre"
	goodUseCase "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/usecase"

	//unitOfMeasurePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure"
	unitOfMeasureRepo "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure/repository/postgre"
	unitOfMeasureUseCase "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure/usecase"
	storePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/storage"
)

type DB struct {
	Timeout time.Duration
	Conn    *pgx.Conn
}

func New(ctx context.Context) storePkg.Interface {
	timeout := time.Duration(time.Millisecond * 8000)
	psqlConn := config.GetConnectionString()
	conn, err := pgx.Connect(ctx, psqlConn)
	if err != nil {
		panic("Unable to connect to database: %v\n" + err.Error())
	}
	if err := conn.Ping(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Connected")

	return &DB{
		Timeout: timeout,
		Conn:    conn,
	}
}

func (c *DB) GetCore(ctx context.Context) *storePkg.Core {
	conn := c.Conn
	timeout := c.Timeout
	rpGood := goodRepo.New(conn, timeout)
	rpCountry := countryRepo.New(conn, timeout)
	rpUnitUfMeasure := unitOfMeasureRepo.New(conn, timeout)
	return &storePkg.Core{
		Good:          goodUseCase.New(rpGood, rpUnitUfMeasure, rpCountry),
		Country:       countryUseCase.New(rpCountry),
		UnitOfMeasure: unitOfMeasureUseCase.New(rpUnitUfMeasure),
	}
}

func (c *DB) Close(ctx context.Context) error {
	defer fmt.Println("Disconnected")
	return c.Conn.Close(ctx)
}
