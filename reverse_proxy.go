package httpx

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func ReverseProxy(addr string) *httputil.ReverseProxy {
	urlObj, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}

	p := httputil.NewSingleHostReverseProxy(urlObj)

	p.Director = func(req *http.Request) {
		req.URL.Scheme = urlObj.Scheme
		req.URL.Host = urlObj.Host

		if req.Header.Get("Authorization") == "" {
			q := req.URL.Query()
			auth := q.Get("auth")
			if auth != "" {
				q.Del("auth")
				req.URL.RawQuery = q.Encode()
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth))
			}
		}
	}

	p.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: time.Hour,
		}).DialContext,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	p.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		slog.Error("cannot proxy request", "trace", r.Header.Get(TraceIDHeader), "error", err)
		http.Error(rw, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}

	p.ModifyResponse = func(response *http.Response) error {
		response.Header.Del("X-Trace-Id")
		response.Header.Del("Access-Control-Expose-Headers")
		return nil
	}

	return p
}
