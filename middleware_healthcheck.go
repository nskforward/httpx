package httpx

import (
	"net/http"
)

func HealthcheckMiddleware(path string) Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) error {
			if ctx.Method() == http.MethodGet && ctx.Path() == path {
				return ctx.RespondText(200, "ok")
			}
			return next(ctx)
		}
	}
}
