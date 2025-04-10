package cors

import (
	"net/http"

	"github.com/nskforward/httpx"
)

func sendPreflight(cfg Config, origin, maxAge string, req *http.Request, resp *httpx.Response) (bool, error) {

	// not a CORS request
	if req.Method != http.MethodOptions {
		return false, nil
	}

	accessMethod := req.Header.Get("Access-Control-Request-Method")
	accessHeaders := req.Header.Get("Access-Control-Request-Headers")

	// not a CORS preflight request
	if accessMethod == "" && accessHeaders == "" {
		return false, nil
	}

	if accessMethod != "" {
		err := sendAllowMethods(cfg, accessMethod, resp)
		if err != nil {
			return false, err
		}
	}

	if accessHeaders != "" {
		err := sendAllowHeaders(cfg, accessHeaders, resp)
		if err != nil {
			return false, err
		}
	}

	sendAllowOrigin(resp, origin)
	sendAllowCredentials(cfg, resp)
	sendExposeHeaders(cfg, resp)
	sendMaxAge(maxAge, resp)

	return true, nil
}

/*
Origin: http://foo.example
Access-Control-Request-Method: POST
Access-Control-Request-Headers: X-PINGOTHER, Content-Type
*/
