package local

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

type cache struct {
	data     map[uint64]models.Good
	lastCode uint64
	mu       sync.RWMutex
	poolCh   chan struct{}
	timeout  time.Duration
}

const poolSize = 10

func New() cachePkg.Interface {
	tm := time.Duration(time.Millisecond * 8000)

	return &cache{data: map[uint64]models.Good{}, mu: sync.RWMutex{}, poolCh: make(chan struct{}, poolSize), timeout: tm}
}

func (c *cache) GetNextCode() uint64 {
	c.lastCode++
	return c.lastCode
}

func (c *cache) List(ctx context.Context) ([]models.Good, error) {

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)

	var result []models.Good

	go func() {
		c.poolCh <- struct{}{}
		c.mu.RLock()
		defer func() {
			c.mu.RUnlock()
			<-c.poolCh
		}()
		result = make([]models.Good, 0, len(c.data))
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

func (c *cache) Add(ctx context.Context, g models.Good) error {
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

func (c *cache) Get(ctx context.Context, code uint64) (*models.Good, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	okChan := make(chan struct{}, 1)
	var ok bool
	var result models.Good
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
			return &result, nil
		}
		return nil, errors.Wrapf(cachePkg.ErrUserNotExists, "code %d", code)
	}

}

func (c *cache) Update(ctx context.Context, g models.Good) error {
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
			err = errors.Wrapf(cachePkg.ErrUserNotExists, "code %d", g.GetCode())
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

func (c *cache) Delete(ctx context.Context, code uint64) error {
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
			err = errors.Wrapf(cachePkg.ErrUserNotExists, "code %d", code)
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

func (c *cache) queueLen() string {
	return fmt.Sprintf(" queue len: %d / %d", len(c.poolCh), poolSize)
}

func (c *cache) Lock() string {
	select {
	case c.poolCh <- struct{}{}:
	default:
		return "queue is full"
	}
	if c.mu.TryLock() {
		return "Write lock: ok" + c.queueLen()
	}
	return "Can't Write lock" + c.queueLen()
}

func (c *cache) RLock() string {
	select {
	case c.poolCh <- struct{}{}:
	default:
		return "queue is full"
	}
	if c.mu.TryRLock() {
		return "Read lock:  ok" + c.queueLen()
	}
	return "Can't Read lock" + c.queueLen()
}

func (c *cache) Unlock() (result string) {
	/* не работает, все равно валится с fatal error
	defer func() {
		if err := recover(); err != nil {
			result = "Write unlock: ERR " + c.queueLen()
		}
	}()
	*/
	select {
	case <-c.poolCh:
	default:
		return "queue is empty"
	}
	c.mu.Unlock()
	result = "Write unlock: ok" + c.queueLen()
	return result
}

func (c *cache) RUnlock() (result string) {
	select {
	case <-c.poolCh:
	default:
		return "queue is empty"
	}
	c.mu.RUnlock()
	return "Read unlock: ok" + c.queueLen()
}
