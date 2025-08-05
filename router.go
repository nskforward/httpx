package httpx

import (
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/nskforward/httpx/mux"
)

type Router struct {
	multiplexer  *mux.Multiplexer
	middlewares  []Middleware
	errorHandler ErrorHandler
}

func NewRouter() *Router {
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
	h, errCode := ro.multiplexer.Search(w, r)
	if errCode > 0 {
		finalHandler := func(w *Response, r *http.Request) error {
			http.Error(w, http.StatusText(errCode), errCode)
			return nil
		}
		for _, mw := range slices.Backward(ro.middlewares) {
			finalHandler = mw(finalHandler)
		}
		resp := newResponse(w)
		err := finalHandler(resp, r)
		if err != nil {
			ro.handlerError(resp, r, err)
		}
		return
	}
	h.ServeHTTP(w, r)
}

func (ro *Router) handlerError(w *Response, r *http.Request, err error) {
	if ro.errorHandler != nil {
		ro.errorHandler(w, r, err)
		return
	}
	fmt.Fprintln(os.Stderr, "httpx: unhandler error during the route", r.Method, r.URL.Path, ">", err)
	if !w.HeadersSent() {
		w.SendText(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}
