package local

import (
	"context"
	"errors"
	"sync"
	"time"

	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

var errNe = errors.New("not implemented")

type cache struct {
	data     map[uint64]*models.Good
	lastCode uint64
	mu       sync.RWMutex
	poolCh   chan struct{}
	timeout  time.Duration
}

const poolSize = 10

func New() cachePkg.Interface {
	tm := time.Duration(time.Millisecond * 8000)

	return &cache{data: map[uint64]*models.Good{}, mu: sync.RWMutex{}, poolCh: make(chan struct{}, poolSize), timeout: tm}
}

func (c *cache) GetNextCode() uint64 {
	c.lastCode++
	return c.lastCode
}

func (c *cache) Close(ctx context.Context) error {
	return nil
}
