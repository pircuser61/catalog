package postgre

import (
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	repoIface "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/unit_of_measure"
)

type countryRepoFixtures struct {
	unitRepo repoIface.Repository
	ctrl     *gomock.Controller
	pool     *pgxpoolmock.MockPgxPool
}

func setUp(t *testing.T) countryRepoFixtures {
	var fixture countryRepoFixtures
	ctrl := gomock.NewController(t)
	mock := pgxpoolmock.NewMockPgxPool(ctrl)
	fixture.pool = mock
	fixture.ctrl = ctrl
	fixture.unitRepo = New(mock, time.Second)
	return fixture
}

func (f *countryRepoFixtures) tearDown() {
	f.ctrl.Finish()
}
