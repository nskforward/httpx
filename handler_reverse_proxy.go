package httpx

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func ReverseProxy(proxyAddr string, rewrite func(r *httputil.ProxyRequest)) Handler {
	rpURL, err := url.Parse(proxyAddr)
	if err != nil {
		panic(err)
	}

	rp := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(rpURL)
			r.SetXForwarded()
			if rewrite != nil {
				rewrite(r)
			}
		},
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: time.Hour,
			}).DialContext,
			TLSHandshakeTimeout:   15 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
		},
	}

	return func(c *Ctx) error {
		rp.ServeHTTP(c.w, c.r)
		return nil
	}
}
