package httpx

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []Handler
	logger      *slog.Logger
}

func NewRouter(logger *slog.Logger) *Router {
	if logger == nil {
		panic("httpx.NewRouter requeres not nil logger")
	}
	return &Router{
		mux:         http.NewServeMux(),
		middlewares: make([]Handler, 0, 16),
		logger:      logger,
	}
}

func (router *Router) Use(middleware ...Handler) {
	router.middlewares = append(router.middlewares, middleware...)
}

func (router *Router) Group(pattern string, middleware ...Handler) *Group {
	return NewGroup(router, pattern, middleware)
}

func (router *Router) DELETE(pattern string, handler Handler, middleware ...Handler) {
	router.CustomMethod("DELETE", pattern, handler, middleware...)
}

func (router *Router) PATCH(pattern string, handler Handler, middleware ...Handler) {
	router.CustomMethod("PATCH", pattern, handler, middleware...)
}

func (router *Router) PUT(pattern string, handler Handler, middleware ...Handler) {
	router.CustomMethod("PUT", pattern, handler, middleware...)
}

func (router *Router) POST(pattern string, handler Handler, middleware ...Handler) {
	router.CustomMethod("POST", pattern, handler, middleware...)
}

func (router *Router) GET(pattern string, handler Handler, middleware ...Handler) {
	router.CustomMethod("GET", pattern, handler, middleware...)
}

func (router *Router) ANY(pattern string, handler Handler, middleware ...Handler) {
	router.CustomMethod("", pattern, handler, middleware...)
}

func (router *Router) CustomMethod(method, pattern string, handler Handler, middleware ...Handler) {
	if handler == nil {
		panic(fmt.Errorf("handler cannot be nil for pattern '%s'", pattern))
	}
	finalHandler := joinHandlers(handler, router.middlewares, middleware)
	if method != "" {
		pattern = fmt.Sprintf("%s %s", method, pattern)
	}
	router.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		resp := NewResponse(router.logger, w, finalHandler)
		err := resp.Next(r)
		router.handleError(resp, err)
	})
}

func (router *Router) handleError(resp *Response, err error) {
	if err != nil {
		apiErr, ok := err.(*APIError)
		if ok {
			resp.Text(apiErr.Code, apiErr.Mesage)
			return
		}
		resp.InternalServerError(err)
		return
	}
	if resp.StatusCode() == 0 {
		resp.InternalServerError(fmt.Errorf("final handler not found"))
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, pattern := router.mux.Handler(r)
	if pattern != "" {
		h.ServeHTTP(w, r)
		return
	}
	router.executeBadHandler(w, r, func(req *http.Request, resp *Response) error {
		h.ServeHTTP(resp.w, req)
		return nil
	})
}

func (router *Router) executeBadHandler(w http.ResponseWriter, r *http.Request, h Handler) {
	resp := NewResponse(router.logger, w, joinHandlers(h, router.middlewares))
	err := resp.Next(r)
	router.handleError(resp, err)
}
