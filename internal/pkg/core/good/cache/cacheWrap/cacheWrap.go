package local

import (
	"context"
	"time"

	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

type cacheWrap struct {
	cache   *cachePkg.Interface
	timeout time.Duration
}

// не понятно по заданию нужно ли принимать на вход котекст, уже с таймаутом
// или context.Background()
// или timeout
func New(_ context.Context, c cachePkg.Interface) cachePkg.Interface {
	tm := time.Duration(time.Millisecond * 8000)

	return &cacheWrap{cache: &c, timeout: tm}
}

func (c *cacheWrap) List() ([]models.Good, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)

	var list []models.Good
	go func() { list, _ = (*c.cache).List(); okChan <- struct{}{} }()
	select {
	case <-ctx.Done():
		return nil, cachePkg.ErrTimeout
	case <-okChan:
		return list, nil
	}

}

func (c *cacheWrap) Add(g models.Good) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)
	var err error
	go func() { err = (*c.cache).Add(g); okChan <- struct{}{} }()
	select {
	case <-ctx.Done():
		return cachePkg.ErrTimeout
	case <-okChan:
		return err
	}
}

func (c *cacheWrap) Get(code uint64) (*models.Good, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)
	var err error
	var g *models.Good
	go func() { g, err = (*c.cache).Get(code); okChan <- struct{}{} }()
	select {
	case <-ctx.Done():
		return nil, cachePkg.ErrTimeout
	case <-okChan:
		return g, err
	}
}

func (c *cacheWrap) Update(g models.Good) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)
	var err error
	go func() { err = (*c.cache).Update(g); okChan <- struct{}{} }()
	select {
	case <-ctx.Done():
		return cachePkg.ErrTimeout
	case <-okChan:
		return err
	}
}

func (c *cacheWrap) Delete(code uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)
	var err error
	go func() { err = (*c.cache).Delete(code); okChan <- struct{}{} }()
	select {
	case <-ctx.Done():
		return cachePkg.ErrTimeout
	case <-okChan:
		return err
	}
}

func (c *cacheWrap) Lock() string {
	return (*c.cache).Lock()
}

func (c *cacheWrap) RLock() string {
	return (*c.cache).RLock()
}

func (c *cacheWrap) Unlock() (result string) {
	return (*c.cache).Unlock()
}

func (c *cacheWrap) RUnlock() (result string) {
	return (*c.cache).RUnlock()
}
