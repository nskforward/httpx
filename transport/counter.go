package transport

import (
	"sync/atomic"

	"golang.org/x/time/rate"
)

type Counter struct {
	value   int64
	limiter *rate.Limiter
}

func (c *Counter) Inc() int64 {
	return atomic.AddInt64(&c.value, 1)
}

func (c *Counter) Dec() int64 {
	return atomic.AddInt64(&c.value, -1)
}
