package counters

import (
	"expvar"
)

var RequestCount, SuccessCount, ErrorCount expvar.Int

func Init() {
	RequestCount.Set(0)
	SuccessCount.Set(0)
	ErrorCount.Set(0)
}

func Request() { RequestCount.Add(1) }

func Success() { SuccessCount.Add(1) }

func Error() { ErrorCount.Add(1) }
