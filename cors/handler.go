package cors

import (
	"net/http"

	"github.com/nskforward/httpx"
)

func CORS(cfg Config) httpx.Handler {

	maxAge := NormalizeMaxAge(cfg)

	return func(req *http.Request, resp *httpx.Response) error {

		origin := req.Header.Get("Origin")
		if origin == "" {
			return resp.Next(req)
		}

		ok := sendPreflight(cfg, origin, maxAge, req, resp)
		if ok {
			return resp.NoContent()
		}

		sendAllowOrigin(cfg, origin, resp)
		sendMaxAge(maxAge, resp)

		return resp.Next(req)
	}
}
