package httpx

import (
	"fmt"
	"time"
)

func CacheDisable(ctx *Context) {
	ctx.SetResponseHeader("Cache-Control", "no-store")
}

func CacheEnable(ctx *Context, public bool, maxAge time.Duration) {
	ctx.SetResponseHeader("Cache-Control", fmt.Sprintf("%s, max-age=%.0f", isPublic(public), maxAge.Seconds()))
}

func isPublic(yes bool) string {
	if yes {
		return "public"
	}
	return "private"
}
