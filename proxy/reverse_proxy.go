package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func ReverseProxy(prefixURL string) *httputil.ReverseProxy {
	addr, err := url.Parse(prefixURL)
	if err != nil {
		panic(err)
	}

	rp := httputil.NewSingleHostReverseProxy(addr)

	rp.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: time.Hour,
		}).DialContext,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
	}

	return rp
}
