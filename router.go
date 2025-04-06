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
	finalHandler := append(router.middlewares, append(middleware, handler)...)
	if method != "" {
		pattern = fmt.Sprintf("%s %s", method, pattern)
	}
	router.mux.HandleFunc(pattern, router.handlerFunc(finalHandler))
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, pattern := router.mux.Handler(r)
	if pattern == "" {
		handler := append(router.middlewares, func(req *http.Request, resp *Response) error {
			h.ServeHTTP(resp.w, req)
			return nil
		})
		router.executeRoute(w, r, handler)
		return
	}
	h.ServeHTTP(w, r)
}

func (router *Router) executeRoute(w http.ResponseWriter, req *http.Request, h []Handler) {
	resp := NewResponse(router.logger, w, req, h)

	err := resp.Next()
	if err != nil {
		fmt.Println("unexpected error:", err)
		resp.InternalServerError()
	}
	if resp.StatusCode() == 0 {
		fmt.Println("unexpeted error:", "request not handled")
		resp.Text(404, "handler not found")
	}
}

func (router *Router) handlerFunc(handler []Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.executeRoute(w, r, handler)
	})
}
