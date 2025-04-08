package cors

import (
	"net/http"

	"github.com/nskforward/httpx"
)

func CORS(cfg Config) httpx.Handler {
	return func(req *http.Request, resp *httpx.Response) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			return resp.Next(req)
		}

	}
}
