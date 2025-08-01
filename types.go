package httpx

import "net/http"

type Multiplexer interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Handle(pattern string, handler http.Handler)
	OnError(h func(w http.ResponseWriter, r *http.Request, code int))
}

type HandlerFunc func(w *Response, r *http.Request) error

type Middleware func(next HandlerFunc) HandlerFunc

type ErrorHandler func(w *Response, r *http.Request, err error)
