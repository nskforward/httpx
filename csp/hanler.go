package csp

import (
	"net/http"

	"github.com/nskforward/httpx"
)

func CSP(cfg Config) httpx.Handler {

	params := Encode(cfg)

	return func(req *http.Request, resp *httpx.Response) error {
		resp.SetHeader("Content-Security-Policy", params)
		return resp.Next(req)
	}
}
