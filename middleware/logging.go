package middleware

import (
	"net/http"

	"github.com/nskforward/httpx/types"
)

func Logging(f func(*types.ResponseWrapper, *http.Request)) types.Middleware {
	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			ww := types.NewResponseWrapper(w)
			err := next(w, r)
			f(ww, r)
			return err
		}
	}
}
