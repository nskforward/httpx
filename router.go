package httpx

import (
	"net/http"

	"github.com/nskforward/httpx/mux"
)

type Router struct {
	multiplexer  *mux.Multiplexer
	middlewares  []Middleware
	errorHandler ErrorHandler
}

func NewRouter() *Router {
	ro := &Router{
		multiplexer:  mux.NewMultiplexer(),
		middlewares:  make([]Middleware, 0, 8),
		errorHandler: defaultErrorHandler,
	}
	return ro
}

func (ro *Router) ErrorHandler(h ErrorHandler) {
	ro.errorHandler = h
}

func (ro *Router) Group(mws ...Middleware) *Router {
	return &Router{
		multiplexer:  ro.multiplexer,
		middlewares:  append(ro.middlewares, mws...),
		errorHandler: ro.errorHandler,
	}
}

func (ro *Router) Use(m ...Middleware) {
	ro.middlewares = append(ro.middlewares, m...)
}

func (ro *Router) Mux() *mux.Multiplexer {
	return ro.multiplexer
}

func (ro *Router) HandleFunc(pattern string, handler HandlerFunc, mws ...Middleware) {
	ro.multiplexer.HandleFunc(pattern, castHandler(ro, handler, mws))
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ro.multiplexer.ServeHTTP(w, r)
}
