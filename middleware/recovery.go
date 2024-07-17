package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/nskforward/httpx/types"
)

func Recovery(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		defer func() {
			if msg := recover(); msg != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				buf = append(buf, '.', '.', '.')
				err = &types.Error{
					Status:     500,
					Text:       fmt.Sprintf("%v", msg),
					StackTrace: string(buf),
				}
			}
		}()
		return next(w, r)
	}
}
