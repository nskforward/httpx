package httpx

import (
	"net"
	"strings"
)

func RealIP(headerName string) Handler {
	return func(ctx *Ctx) error {
		ip := ctx.Request().Header.Get(headerName)
		if ip == "" {
			host, _, _ := net.SplitHostPort(ctx.r.RemoteAddr)
			if host != "" {
				ctx.r.RemoteAddr = host
			}
			return ctx.Next()
		}
		i := strings.Index(ip, ",")
		if i > -1 {
			ctx.r.RemoteAddr = strings.TrimSpace(ip[:i])
		} else {
			ctx.r.RemoteAddr = ip
		}
		return ctx.Next()
	}
}
