package middleware

import (
	"net/http"
	"strings"

	"github.com/nskforward/httpx/types"
)

func RealIP(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ip := proxiedIP(r.Header)
		if ip != "" {
			r.RemoteAddr = ip
		}
		return next(w, r)
	}
}

func proxiedIP(header http.Header) string {
	ip := header.Get(types.XForwardedFor)
	if ip != "" {
		i := strings.LastIndex(ip, ",")
		if i > -1 {
			ip = strings.TrimSpace(ip[i+1:])
			return ip
		}
	}

	return header.Get(types.XRealIP)
}
