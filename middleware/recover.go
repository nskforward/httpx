package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/nskforward/httpx/types"
)

func Recover(onPanic func(err error, trace string)) types.Middleware {
	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			defer func() {
				if err := recover(); r != nil {
					if err == http.ErrAbortHandler {
						panic(err)
					}
					onPanic(fmt.Errorf("%s", err), string(debug.Stack()))
					http.Error(w, "Internal Server Error", 500)
				}
			}()
			return next(w, r)
		}
	}
}
