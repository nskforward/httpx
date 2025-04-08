package hsts

import (
	"net/http"

	"github.com/nskforward/httpx"
)

func HSTS(cfg Config) httpx.Handler {

	value := cfg.Encode()

	return func(req *http.Request, resp *httpx.Response) error {
		resp.SetHeader("Strict-Transport-Security", value)
		return resp.Next(req)
	}
}
