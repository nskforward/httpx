package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
	"golang.org/x/time/rate"
)

var rlStore sync.Map

func RateLimiter(period time.Duration, burst int, timeout time.Duration, keyFunc func(r *http.Request) string) types.Middleware {
	if keyFunc == nil {
		panic("rate limiter keyFunc function cannot be a nil")
	}

	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			key := keyFunc(r)
			rl, ok := rlStore.Load(key)
			if !ok {
				rl = rate.NewLimiter(rate.Every(period), burst)
				rlStore.Store(key, rl)
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			err := rl.(*rate.Limiter).Wait(ctx)
			if err != nil {
				return response.Error{Status: http.StatusTooManyRequests}
			}
			return next(w, r)
		}
	}
}
