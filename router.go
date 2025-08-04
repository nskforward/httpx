package httpx

import (
	"log/slog"
	"net/http"

	"github.com/nskforward/httpx/mux"
)

type Router struct {
	logger       *slog.Logger
	multiplexer  *mux.Multiplexer
	middlewares  []Middleware
	errorHandler ErrorHandler
}

func NewRouter(logger *slog.Logger) *Router {
	return &Router{
		logger:      logger,
		multiplexer: mux.NewMultiplexer(),
		middlewares: make([]Middleware, 0, 8),
	}
}

func (ro *Router) ErrorHandler(h ErrorHandler) {
	ro.errorHandler = h
}

func (ro *Router) Group(mws ...Middleware) *Router {
	return &Router{
		logger:       ro.logger,
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

func (ro *Router) handlerError(w *Response, r *http.Request, err error) {
	if ro.errorHandler != nil {
		ro.errorHandler(w, r, err)
		return
	}
	ro.logger.Error("httpx: unhandler error during the route", "method", r.Method, "path", r.URL.Path, "error", err)
	if !w.HeadersSent() {
		w.SendShortError(http.StatusInternalServerError)
	}
}
