package httpx

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/nskforward/httpx/mux"
)

type Router struct {
	multiplexer  *mux.Multiplexer
	middlewares  []Middleware
	errorHandler ErrorHandler
}

func NewRouter(logger *slog.Logger) *Router {
	return &Router{
		multiplexer: mux.NewMultiplexer(),
		middlewares: make([]Middleware, 0, 8),
	}
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

func (ro *Router) handlerError(w *Response, r *http.Request, err error) {
	if ro.errorHandler != nil {
		ro.errorHandler(w, r, err)
		return
	}
	fmt.Fprintln(os.Stderr, "httpx: unhandler error during the route", r.Method, r.URL.Path, ">", err)
	if !w.HeadersSent() {
		w.SendShortError(http.StatusInternalServerError)
	}
}
