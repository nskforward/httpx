package httpx

import (
	"net/http"
	"strings"
)

func TrailingSlash(ctx *Ctx) error {
	if ctx.Path() == "/" || !strings.HasSuffix(ctx.Path(), "/") {
		return ctx.Next()
	}

	u := ctx.Request().URL
	u.Path = strings.TrimRight(u.Path, "/")

	return ctx.Redirect(http.StatusPermanentRedirect, u.String())
}
