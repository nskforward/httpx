package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nskforward/httpx/types"
)

func RealIP(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ip := proxiedIP(r.Header)
		if ip != "" {
			fmt.Println("RemoteAddr:", ip)
			r.RemoteAddr = ip
		} else {
			fmt.Println("ip not changed:", ip)
		}
		return next(w, r)
	}
}

func proxiedIP(header http.Header) string {
	ip := header.Get(types.XForwardedFor)
	if ip != "" {
		fmt.Println("XForwardedFor:", ip)
		i := strings.Index(ip, ",")
		if i > -1 {
			ip = strings.TrimSpace(ip[:i])
			return ip
		}
	}

	return header.Get(types.XRealIP)
}
