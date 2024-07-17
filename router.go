package httpx

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"

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
	router.mux.HandleFunc(pattern, finalHandler(h, router.middlewares, middlewares))
}

func (router *Router) Group(middleware ...types.Middleware) *Router {
	return &Router{
		mux:         router.mux,
		middlewares: append(router.middlewares, middleware...),
	}
}

func (router *Router) Use(middlewares ...types.Middleware) {
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
		if err != nil {
			resp, ok := err.(types.Error)
			if ok {
				slog.Error(fmt.Sprintf("http-%d", resp.Status), "error", resp.Text, "trace-id", r.Header.Get(types.XTraceID), "stacktrace", resp.StackTrace)
				if resp.Text == "" || resp.Status == 500 {
					http.Error(w, http.StatusText(resp.Status), resp.Status)
				} else {
					http.Error(w, resp.Text, resp.Status)
				}
			} else {
				slog.Error("http-400", "error", err.Error(), "trace-id", r.Header.Get(types.XTraceID))
				http.Error(w, err.Error(), 400)
			}
		}
	}
}
