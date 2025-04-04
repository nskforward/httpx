package httpx

import (
	"net/http"
	"strings"
)

type router struct {
	mux            *http.ServeMux
	app            *App
	appMiddlewares []Handler
}

func newRouter(app *App) *router {
	return &router{
		mux: http.NewServeMux(),
		app: app,
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

func (r *router) Route(method, pattern string, handler Handler, middlewares []Handler) {
	if strings.Contains(pattern, " ") {
		panic("white spaces not allowed in http router patterns")
	}
	handlers := append(r.appMiddlewares, middlewares...)
	route := newRoute(r, pattern, append(handlers, handler))
	route.registry(method)
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
