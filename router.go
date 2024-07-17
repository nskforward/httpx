package httpx

import (
	"crypto/tls"
	"log/slog"
	"net/http"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/transport"
	"github.com/nskforward/httpx/types"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []types.Middleware
}

/*
Patterns can match the method, host and path of a request. Some examples:

	"/index.html" matches the path "/index.html" for any host and method.
	"GET /static/" matches a GET request whose path begins with "/static/".
	"example.com/" matches any request to the host "example.com".
	"example.com/{$}" matches requests with host "example.com" and path "/".
	"/b/{bucket}/o/{objectname...}" matches paths whose first segment is "b" and whose third segment is "o".
*/
func (router *Router) Route(pattern string, h types.Handler, middlewares ...types.Middleware) {
	if router.mux == nil {
		router.mux = http.NewServeMux()
	}
	if router.middlewares == nil {
		router.middlewares = make([]types.Middleware, 0, 8)
	}
	router.mux.HandleFunc(pattern, finalHandler(h, router.middlewares, middlewares))
}

func (router *Router) Group(middleware ...types.Middleware) *Router {
	if router.mux == nil {
		router.mux = http.NewServeMux()
	}
	if router.middlewares == nil {
		router.middlewares = make([]types.Middleware, 0, 8)
	}
	return &Router{
		mux:         router.mux,
		middlewares: append(router.middlewares, middleware...),
	}
}

func (router *Router) Use(middlewares ...types.Middleware) {
	if router.middlewares == nil {
		router.middlewares = make([]types.Middleware, 0, 8)
	}
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *Router) Listen(addr string) error {
	return transport.DefaultTransport().Listen(addr, router)
}

func (router *Router) ListenTLS(addr string, tlsConfig *tls.Config) error {
	return transport.DefaultTransport().ListenTLS(addr, tlsConfig, router)
}

func finalHandler(h types.Handler, mw1, mw2 []types.Middleware) http.HandlerFunc {
	if mw1 == nil && mw2 == nil {
		return nil
	}
	mw0 := make([]types.Middleware, 0, len(mw1)+len(mw2))
	mw0 = append(mw0, mw1...)
	mw0 = append(mw0, mw2...)
	for i := len(mw0) - 1; i >= 0; i-- {
		h = mw0[i](h)
	}
	return catch(h)
}

func catch(next types.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err == nil {
			return
		}
		resp, ok := err.(response.Error)
		if ok {
			slog.Error("unhandled error", "status", resp.Status, "error", resp.Text)
			http.Error(w, resp.Text, resp.Status)
			return
		}
		slog.Error("unhandled error", "status", resp.Status, "error", err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

/*
HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Type: text/html; charset=utf-8
Date: Sun, 16 Jun 2024 09:56:18 GMT
ETag: W/"34d80-RRdtGW1ieWo+bM5DRkLbeL33cQQ"
Server: Tank
Strict-Transport-Security: max-age=31536000; includeSubDomains
Transfer-Encoding: Identity
Vary: Accept-Encoding, Accept-Encoding
X-Content-Type-Options: nosniff
X-DNS-Prefetch-Control: off
X-Download-Options: noopen
X-Frame-Options: SAMEORIGIN
X-Request-Detected-Device: desktop
X-Request-Geoip-Country-Code: RU
X-Request-Id: fd447a85a55d4d16944b656fc2cea2de
X-XSS-Protection: 1; mode=block

HTTP/1.1 200 OK
Content-Type: text/plain
Transfer-Encoding: chunked

0\r\n
Mozilla\r\n
7\r\n
Developer\r\n
9\r\n
Network\r\n
0\r\n
\r\n
*/
