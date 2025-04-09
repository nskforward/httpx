package cors

import (
	"net/http"

	"github.com/nskforward/httpx"
)

func CORS(cfg Config) httpx.Handler {

	cfg.AllowOrigins = normalizeOrigins(cfg.AllowOrigins)
	maxAge := NormalizeMaxAge(cfg)

	return func(req *http.Request, resp *httpx.Response) error {

		origin := req.Header.Get("Origin")
		if origin == "" {
			return resp.Next(req)
		}

		ok, err := sendPreflight(cfg, origin, maxAge, req, resp)
		if err != nil {
			resp.Logger().Warn("CORS validation failed", "error", err.Error())
			return resp.Forbidden(err.Error())
		}
		if ok {
			return resp.NoContent()
		}

		err = sendAllowOrigin(cfg, origin, resp)
		if err != nil {
			resp.Logger().Warn("CORS validation failed", "error", err.Error())
			return resp.Forbidden(err.Error())
		}

		sendMaxAge(maxAge, resp)

		return resp.Next(req)
	}
}
