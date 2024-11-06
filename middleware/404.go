package middleware

import (
	"net/http"

	"github.com/nskforward/httpx/types"
)

func NotFound(handler func(http.ResponseWriter, *http.Request)) types.Middleware {
	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			ww := types.NewResponseWrapper(w)
			ww.BeforeBody = func() {
				if ww.Status() == 404 {
					handler(ww, r)
					ww.SkipBody()
				}
			}
			return next(w, r)
		}
	}
}
