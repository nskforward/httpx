package httpx

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/nskforward/httpx/cache"
	"github.com/nskforward/httpx/transport"
	"github.com/nskforward/httpx/types"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []types.Middleware
	loggerFunc  types.LoggerFunc
	errorFunc   types.ErrorFunc
}

func NewRouter() *Router {
	r := &Router{
		mux:         http.NewServeMux(),
		middlewares: make([]types.Middleware, 0, 8),
		loggerFunc:  defaultLoggerFunc,
		errorFunc:   defaultErrorFunc,
	}
	return r
}

func (router *Router) LoggerFunc(f types.LoggerFunc) {
	router.loggerFunc = f
}

func (router *Router) ErrorFunc(f types.ErrorFunc) {
	router.errorFunc = f
}

/*
Patterns can match the method, host and path of a request. Some examples:

	"/index.html" matches the path "/index.html" for any host and method.
	"GET /static/" matches a GET request whose path begins with "/static/".
	"example.com/" matches any request to the host "example.com".
	"example.com/{$}" matches requests with host "example.com" and path "/".
	"/b/{bucket}/o/{objectname...}" matches paths whose first segment is "b" and whose third segment is "o".
*/
func (router *Router) Route(pattern string, h types.Handler, middlewares ...types.Middleware) *Router {
	if router.mux == nil {
		panic(fmt.Errorf("uninitialized router"))
	}
	router.mux.HandleFunc(pattern, router.handler(h, middlewares))
	return router
}

func (router *Router) RouteH(pattern string, h http.Handler, middlewares ...types.Middleware) *Router {
	return router.Route(
		pattern,
		func(w http.ResponseWriter, r *http.Request) error {
			h.ServeHTTP(w, r)
			return nil
		},
		middlewares...,
	)
}

func (router *Router) RouteHF(pattern string, h http.HandlerFunc, middlewares ...types.Middleware) *Router {
	return router.Route(
		pattern,
		func(w http.ResponseWriter, r *http.Request) error {
			h.ServeHTTP(w, r)
			return nil
		},
		middlewares...,
	)
}

func (router *Router) RouteV(pattern string, comp templ.Component, middlewares ...types.Middleware) *Router {
	return router.Route(
		pattern,
		func(w http.ResponseWriter, r *http.Request) error {
			cache.Prohibit(w)
			templ.Handler(comp).ServeHTTP(w, r)
			return nil
		},
		middlewares...,
	)
}

func (router *Router) Group(middleware ...types.Middleware) *Router {
	if router.mux == nil {
		panic(fmt.Errorf("uninitialized router"))
	}
	return &Router{
		mux:         router.mux,
		middlewares: append(router.middlewares, middleware...),
	}
}

func (router *Router) Use(middlewares ...types.Middleware) {
	if router.mux == nil {
		panic(fmt.Errorf("uninitialized router"))
	}
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if router.mux == nil {
		panic(fmt.Errorf("uninitialized router"))
	}
	h := func(ww http.ResponseWriter, rr *http.Request) error {
		router.mux.ServeHTTP(ww, rr)
		return nil
	}
	for i := len(router.middlewares) - 1; i >= 0; i-- {
		h = router.middlewares[i](h)
	}
	ww := types.NewResponseWrapper(w)
	router.Catch(h)(ww, r)
	router.loggerFunc(ww, r)
}

func (router *Router) Listen(addr string) error {
	return transport.DefaultTransport().Listen(addr, router)
}

func (router *Router) ListenTLS(addr string, tlsConfig *tls.Config) error {
	return transport.DefaultTransport().ListenTLS(addr, tlsConfig, router)
}

func (router *Router) handler(h types.Handler, middleware []types.Middleware) http.HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return router.Catch(h)
}

/*
HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Type: text/html; charset=utf-8
Date: Sun, 16 Jun 2024 09:56:18 GMT
ETag: W/"34d80-RRdtGW1ieWo+bM5DRkLbeL33cQQ"
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
