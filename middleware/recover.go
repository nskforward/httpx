package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/nskforward/httpx/types"
)

func Recover(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer func() {
			if r := recover(); r != nil {
				if r == http.ErrAbortHandler {
					panic(r)
				}
				http.Error(w, "Internal Server Error", 500)
				fmt.Println("FATAL", r, "\n", string(debug.Stack()))
			}
		}()
		return next(w, r)
	}
}
