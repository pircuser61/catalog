package counters

import (
	"expvar"
)

var RequestCount, SuccessCount, ErrorCount, CacheHit, CacheMiss expvar.Int

func Init() {
	RequestCount.Set(0)
	SuccessCount.Set(0)
	ErrorCount.Set(0)
	CacheHit.Set(0)
	CacheMiss.Set(0)
}

func Request() { RequestCount.Add(1) }

func Success() { SuccessCount.Add(1) }

func Error() { ErrorCount.Add(1) }

func Hit() { CacheHit.Add(1) }

func Miss() { CacheMiss.Add(1) }
