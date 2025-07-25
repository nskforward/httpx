package httpx

import "net/http"

type Multiplexer interface {
	HandleFunc(pattern string, handler http.HandlerFunc)
	Handle(pattern string, handler http.Handler)
}
