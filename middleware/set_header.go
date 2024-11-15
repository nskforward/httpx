package middleware

import (
	"net/http"

	"github.com/nskforward/httpx/types"
)

func SetHeaders(headers map[string]string) types.Middleware {
	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			for h, v := range headers {
				w.Header().Add(h, v)
			}
			return next(w, r)
		}
	}
}
