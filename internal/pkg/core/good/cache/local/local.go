package local

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	cachePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/cache"
	"gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good/models"
)

type cache struct {
	data     map[uint64]models.Good
	lastCode uint64
	mu       sync.RWMutex
	poolCh   chan struct{}
}

const poolSize = 4

func New() cachePkg.Interface {
	return &cache{data: map[uint64]models.Good{}, mu: sync.RWMutex{}, poolCh: make(chan struct{}, poolSize)}
}

func (c *cache) GetNextCode() uint64 {
	c.lastCode++
	return c.lastCode
}

func (c *cache) List() []models.Good {
	c.poolCh <- struct{}{}
	c.mu.RLock()
	defer func() {
		c.mu.RUnlock()
		<-c.poolCh
	}()
	result := make([]models.Good, 0, len(c.data))
	for _, x := range c.data {
		result = append(result, x)
	}
	return result
}

func (c *cache) Add(g models.Good) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()
	//if _, ok := c.data[g.GetCode()]; ok {
	//	return errors.Wrapf(ErrUserExists, "code %d", g.GetCode())
	//}
	if err := g.SetCode(c.GetNextCode()); err != nil {
		return err
	}
	c.data[g.GetCode()] = g
	return nil
}

func (c *cache) Get(code uint64) (*models.Good, error) {
	c.poolCh <- struct{}{}
	c.mu.RLock()
	defer func() {
		c.mu.RUnlock()
		<-c.poolCh
	}()
	if g, ok := c.data[code]; ok {
		return &g, nil
	}
	return nil, errors.Wrapf(cachePkg.ErrUserNotExists, "code %d", code)
}

func (c *cache) Update(g models.Good) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()

	if _, ok := c.data[g.GetCode()]; !ok {
		return errors.Wrapf(cachePkg.ErrUserNotExists, "code %d", g.GetCode())
	}
	c.data[g.GetCode()] = g
	return nil
}

func (c *cache) Delete(code uint64) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		<-c.poolCh
	}()
	if _, ok := c.data[code]; ok {
		delete(c.data, code)
		return nil
	}
	return errors.Wrapf(cachePkg.ErrUserNotExists, "code %d", code)
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
