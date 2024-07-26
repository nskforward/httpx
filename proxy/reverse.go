package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/nskforward/httpx/types"
)

func Reverse(prefixURL string) types.Handler {
	rpURL, err := url.Parse(prefixURL)
	if err != nil {
		panic(err)
	}

	rp := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(rpURL)
			r.SetXForwarded()
			r.Out.Header.Del(types.XRealIP)
		},
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: time.Hour,
			}).DialContext,
			TLSHandshakeTimeout:   15 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		rp.ServeHTTP(w, r)
		return nil
	}
}
