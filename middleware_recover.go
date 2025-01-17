package httpx

import (
	"fmt"
	"net/http"
)

func RecoverMiddleware(next Handler) Handler {
	return func(ctx *Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)

				}
				if err == http.ErrAbortHandler {
					panic(err)
				}
				ctx.Logger().Error("panic", "error", err)
				ctx.RespondText(http.StatusInternalServerError, fmt.Sprintf("internal server error: trace id: %s", ctx.TraceID()))
			}
		}()
		return next(ctx)
	}
}
