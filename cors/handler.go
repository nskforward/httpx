package cors

import (
	"net/http"
	"strconv"
	"time"

	"github.com/nskforward/httpx"
)

func CORS(cfg Config) httpx.Handler {

	if len(cfg.AllowOrigins) == 0 {
		panic("cors config AllowOrigins field cannot be empty")
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = time.Minute
	}

	maxAge := strconv.Itoa(int(cfg.MaxAge.Seconds()))

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
