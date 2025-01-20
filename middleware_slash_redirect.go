package httpx

import (
	"net/http"
	"strings"
)

func SlashRedirectMiddleware(next Handler) Handler {
	return func(ctx *Context) error {
		if ctx.Path() == "/" || !strings.HasSuffix(ctx.Path(), "/") {
			return next(ctx)
		}

		u := ctx.Request().URL
		u.Path = strings.TrimRight(u.Path, "/")

		return ctx.Redirect(http.StatusPermanentRedirect, u.String())
	}
}
