package local

import (
	"context"

	"github.com/pkg/errors"
	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

func (c *cache) GoodList(ctx context.Context) ([]*models.Good, error) {

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)

	var result []*models.Good

	go func() {
		c.poolCh <- struct{}{}
		c.mu.RLock()
		defer func() {
			c.mu.RUnlock()
			<-c.poolCh
		}()
		result = make([]*models.Good, 0, len(c.data))
		for _, x := range c.data {
			result = append(result, x)
		}
		okChan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, cachePkg.ErrTimeout
	case <-okChan:
		return result, nil
	}

}

func (c *cache) GoodAdd(ctx context.Context, g *models.Good) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)
	var err error
	go func() {
		c.poolCh <- struct{}{}
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()
		err = g.SetCode(c.GetNextCode())
		if err == nil {
			c.data[g.GetCode()] = g
		}
		okChan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return cachePkg.ErrTimeout
	case <-okChan:
		return err
	}

}

func (c *cache) GoodGet(ctx context.Context, code uint64) (*models.Good, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)
	var ok bool
	var result *models.Good
	go func() {
		c.poolCh <- struct{}{}
		c.mu.RLock()
		defer func() {
			c.mu.RUnlock()
			<-c.poolCh
		}()
		result, ok = c.data[code]
		okChan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, cachePkg.ErrTimeout
	case <-okChan:
		if ok {
			return result, nil
		}
		return nil, errors.Wrapf(cachePkg.ErrObjNotExists, "code %d", code)
	}

}

func (c *cache) GoodUpdate(ctx context.Context, g *models.Good) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)
	var err error
	go func() {
		c.poolCh <- struct{}{}
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()

		if _, ok := c.data[g.GetCode()]; !ok {
			err = errors.Wrapf(cachePkg.ErrObjNotExists, "code %d", g.GetCode())
		} else {
			c.data[g.GetCode()] = g
		}
		okChan <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return cachePkg.ErrTimeout
	case <-okChan:
		return err
	}

}

func (c *cache) GoodDelete(ctx context.Context, code uint64) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)
	var err error
	go func() {
		c.poolCh <- struct{}{}
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()
		if _, ok := c.data[code]; ok {
			delete(c.data, code)
			err = nil
		} else {
			err = errors.Wrapf(cachePkg.ErrObjNotExists, "code %d", code)
		}
		okChan <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return cachePkg.ErrTimeout
	case <-okChan:
		return err
	}
}
