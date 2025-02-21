package httpx

import (
	"context"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type KeyFunc func(r *http.Request) string
type SkipFunc func(r *http.Request) bool

var rlStore sync.Map

func RateLimiter(period, timeout time.Duration, maxHits int, keyFunc KeyFunc, skip SkipFunc) Handler {
	return func(ctx *Ctx) error {
		if skip != nil && skip(ctx.Request()) {
			return ctx.Next()
		}

		key := keyFunc(ctx.Request())
		if key == "" {
			return ctx.Next()
		}

		rl, ok := rlStore.Load(key)
		if !ok {
			rl = rate.NewLimiter(rate.Every(period), maxHits)
			rlStore.Store(key, rl)
			return ctx.Next()
		}

		waitContext, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := rl.(*rate.Limiter).Wait(waitContext)
		if err != nil {
			return ErrTooManyRequests
		}

		return ctx.Next()
	}
}
