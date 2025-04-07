package middleware

import (
	"net/http"
	"strings"

	"github.com/nskforward/httpx"
)

func RealIP(header string) httpx.Handler {
	return func(req *http.Request, resp *httpx.Response) error {
		addr := req.Header.Get(header)
		if addr != "" {
			i := strings.Index(addr, ",")
			if i > -1 {
				addr = strings.TrimSpace(addr[:i])
			}
			req.RemoteAddr = addr
		}
		return resp.Next(req)
	}
}
