package httpx

import (
	"log/slog"
	"net/http"
)

type Router struct {
	serverMux        *http.ServeMux
	logger           *slog.Logger
	mws              []Middleware
	beforeRequestLog LogFunc
	afterResponseLog LogFunc
	slashRedirect    bool
}

func NewRouter(logger *slog.Logger, opts ...SetOpt) *Router {
	if logger == nil {
		panic("logger cannot be nil")
	}

	r := &Router{
		serverMux:     http.NewServeMux(),
		mws:           make([]Middleware, 0, 16),
		slashRedirect: true,
		logger:        logger,
	}
	for _, opt := range opts {
		opt(r)
	}
	if r.slashRedirect {
		r.Use(SlashRedirectMiddleware)
	}
	return r
}

func (r *Router) Use(middleware ...Middleware) {
	r.mws = append(r.mws, middleware...)
}

func (r *Router) Group(patternPrefix string) *Group {
	return NewGroup(r, patternPrefix)
}

func (r *Router) ANY(pattern string, handler Handler, middlewares ...Middleware) {
	DeclareHandler(r, "", pattern, handler, middlewares...)
}

func (r *Router) GET(pattern string, handler Handler, middlewares ...Middleware) {
	DeclareHandler(r, "GET", pattern, handler, middlewares...)
}

func (r *Router) POST(pattern string, handler Handler, middlewares ...Middleware) {
	DeclareHandler(r, "POST", pattern, handler, middlewares...)
}

func (r *Router) PUT(pattern string, handler Handler, middlewares ...Middleware) {
	DeclareHandler(r, "PUT", pattern, handler, middlewares...)
}

func (r *Router) DELETE(pattern string, handler Handler, middlewares ...Middleware) {
	DeclareHandler(r, "DELETE", pattern, handler, middlewares...)
}

func (r *Router) PATCH(pattern string, handler Handler, middlewares ...Middleware) {
	DeclareHandler(r, "PATCH", pattern, handler, middlewares...)
}

func (r *Router) OPTIONS(pattern string, handler Handler, middlewares ...Middleware) {
	DeclareHandler(r, "OPTIONS", pattern, handler, middlewares...)
}

func (r *Router) HEAD(pattern string, handler Handler, middlewares ...Middleware) {
	DeclareHandler(r, "HEAD", pattern, handler, middlewares...)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	stdhandler, pattern := r.serverMux.Handler(req)
	if pattern != "" {
		stdhandler.ServeHTTP(w, req)
		return
	}
	handler := func(ctx *Context) error {
		stdhandler.ServeHTTP(ctx.w, ctx.req)
		return nil
	}
	executeFinalHandler(r, handler, r.mws, w, req)
}
