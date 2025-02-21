package httpx

import (
	"fmt"
	"net/http"
	"strings"
)

type Route struct {
	router   *router
	pattern  string
	handlers []Handler
}

func newRoute(router *router, pattern string, handlers []Handler) *Route {
	return &Route{
		router:   router,
		pattern:  pattern,
		handlers: handlers,
	}
}

func (route *Route) registry(method Method) {
	if strings.Contains(route.pattern, " ") {
		panic(fmt.Errorf("white spaces not allowed in http router pattern: %s", route.pattern))
	}
	pattern := route.pattern
	if method != ANY {
		pattern = fmt.Sprintf("%s %s", method, route.pattern)
	}
	route.router.mux.Handle(pattern, route)
}

func (route *Route) Route(method Method, pattern string, handler Handler, middlewares ...Handler) *Route {
	handlers := append(route.handlers, middlewares...)
	finalPattern := route.pattern + pattern
	r := newRoute(route.router, finalPattern, append(handlers, handler))
	r.registry(method)
	return r
}

func (route *Route) Group(pattern string, middlewares ...Handler) *Route {
	return newRoute(route.router, route.pattern+pattern, append(route.handlers, middlewares...))
}

func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newCtx(route, w, r)
	for {
		if ctx.Sent() {
			break
		}
		nextHandler := ctx.nextHandler()
		if nextHandler == nil {
			break
		}
		err := nextHandler(ctx)
		if err != nil {
			ctx.WriteError(err)
			break
		}
	}
	if !ctx.Sent() {
		ctx.WriteError(ErrNotFound)
	}
}
