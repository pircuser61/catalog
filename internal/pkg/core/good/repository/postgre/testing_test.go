package postgre

import (
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	repoIface "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
)

type goodsRepoFixtures struct {
	goodRepo repoIface.Repository
	ctrl     *gomock.Controller
	pool     *pgxpoolmock.MockPgxPool
}

func setUp(t *testing.T) goodsRepoFixtures {
	var fixture goodsRepoFixtures
	ctrl := gomock.NewController(t)
	mock := pgxpoolmock.NewMockPgxPool(ctrl)
	fixture.pool = mock
	fixture.ctrl = ctrl
	fixture.goodRepo = New(mock, time.Second)
	return fixture
}

func (f *goodsRepoFixtures) tearDown() {
	f.ctrl.Finish()
}
