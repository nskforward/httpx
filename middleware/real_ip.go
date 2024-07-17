package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/nskforward/httpx/types"
)

func SetRealIP(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ip := proxiedIP(r.Header)
		if ip != "" {
			r.RemoteAddr = ip
		}
		return next(w, r)
	}
}

func proxiedIP(header http.Header) string {
	var ip string
	ip = header.Get(types.TrueClientIP)
	if ip == "" {
		ip = header.Get(types.XRealIP)
		if ip == "" {
			ip = header.Get(types.XForwardedFor)
			if ip == "" {
				return ""
			}
			i := strings.Index(ip, ",")
			if i > -1 {
				ip = ip[:i]
			}
		}
	}
	if net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}
