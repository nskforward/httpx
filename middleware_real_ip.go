package httpx

import (
	"net/http"
	"strings"
)

func RealIPMiddleware(next Handler) Handler {
	return func(ctx *Context) error {
		ip := proxiedIP(ctx.Request().Header)
		if ip != "" {
			ctx.Request().RemoteAddr = ip
		}
		return next(ctx)
	}
}

func proxiedIP(header http.Header) string {
	ip := header.Get("X-Forwarded-For")
	if ip != "" {
		i := strings.Index(ip, ",")
		if i > -1 {
			ip = strings.TrimSpace(ip[:i])
		}
		return ip
	}
	return header.Get("X-Real-IP")
}
