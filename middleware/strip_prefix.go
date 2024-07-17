package middleware

import (
	"net/http"
	"strings"

	"github.com/nskforward/httpx/types"
)

func StripPrefix(prefix string) types.Middleware {
	return func(next types.Handler) types.Handler {
		if prefix == "" {
			return next
		}
		return func(w http.ResponseWriter, r *http.Request) error {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, prefix)
			return next(w, r)
		}
	}
}
