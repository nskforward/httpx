package cors

import (
	"net/http"

	"github.com/nskforward/httpx"
)

func sendPreflight(cfg Config, origin, maxAge string, req *http.Request, resp *httpx.Response) bool {
	if req.Method != http.MethodOptions {
		return false
	}

	accessMethod := req.Header.Get("Access-Control-Request-Method")
	accessHeaders := req.Header.Get("Access-Control-Request-Headers")
	if accessMethod == "" && accessHeaders == "" {
		return false
	}

	if accessMethod != "" {
		sendAllowMethods(cfg, resp)
	}
	if accessHeaders != "" {
		sendAllowHeaders(cfg, resp)
	}

	sendAllowOrigin(cfg, origin, resp)
	sendAllowCredentials(cfg, resp)
	sendExposeHeaders(cfg, resp)
	sendMaxAge(maxAge, resp)

	return true
}

/*
Origin: http://foo.example
Access-Control-Request-Method: POST
Access-Control-Request-Headers: X-PINGOTHER, Content-Type
*/
