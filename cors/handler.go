package cors

import (
	"fmt"
	"net/http"

	"github.com/nskforward/httpx"
)

func CORS(cfg Config) httpx.Handler {
	originPool := ParseOriginPool(cfg.AllowLocalhost, cfg.AllowOrigins)
	maxAge := NormalizeMaxAge(cfg)

	return func(req *http.Request, resp *httpx.Response) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			return resp.Next(req)
		}

		if !originPool.Valid(origin) {
			resp.Logger().Warn("cors validation", "error", "origin not allowed", "origin", origin)
			return resp.Forbidden(fmt.Sprintf("cors: origin '%s' not allowed", origin))
		}

		ok, err := sendPreflight(cfg, origin, maxAge, req, resp)
		if err != nil {
			resp.Logger().Warn("cors validation", "error", err.Error())
			return resp.Forbidden(err.Error())
		}
		if ok {
			return resp.NoContent()
		}

		sendAllowOrigin(resp, origin)
		sendMaxAge(maxAge, resp)

		return resp.Next(req)
	}
}
