package middleware

import (
	"net/http"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func SetHeader(name, value string, once bool) types.Middleware {
	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if !once {
				w.Header().Set(name, value)
				return next(w, r)
			}
			ww := response.NewWrapper(w)
			ww.BeforeBody = func() {
				if ww.Header().Get(name) == "" {
					ww.Header().Set(name, value)
				}
			}
			return next(ww, r)
		}
	}
}
