package local

import "fmt"

/*
	Методы для тестирования
*/

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
