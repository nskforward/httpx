package httpx

import (
	"net/http"
	"strings"
)

type router struct {
	mux            *http.ServeMux
	server         *Server
	appMiddlewares []Handler
}

func newRouter(s *Server) *router {
	return &router{
		mux:    http.NewServeMux(),
		server: s,
	}
}

func (r *router) wrapStdHandler(h http.Handler) http.Handler {
	finalHandler := func(c *Ctx) error {
		h.ServeHTTP(c.w, c.r)
		return nil
	}
	return newRoute(r, "", append(r.appMiddlewares, finalHandler))
}

func (r *router) use(middlewares []Handler) {
	if r.appMiddlewares == nil {
		r.appMiddlewares = middlewares
	} else {
		r.appMiddlewares = append(r.appMiddlewares, middlewares...)
	}
}

func (r *router) Route(method Method, pattern string, handler Handler, middlewares []Handler) *Route {
	if strings.Contains(pattern, " ") {
		panic("white spaces not allowed in http router patterns")
	}
	handlers := append(r.appMiddlewares, middlewares...)
	route := newRoute(r, pattern, append(handlers, handler))
	route.registry(method)
	return route
}

func (r *router) Group(pattern string, middlewares []Handler) *Route {
	return newRoute(r, pattern, append(r.appMiddlewares, middlewares...))
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h, pattern := r.mux.Handler(req)

	if pattern != "" {
		h.ServeHTTP(w, req)
		return
	}

	h = r.wrapStdHandler(h)
	h.ServeHTTP(w, req)
}
