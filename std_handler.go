package httpx

import (
	"net/http"
)

func toStdHandler(r *Router, pattern string, handler Handler, mws []Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		executeHandler(r, pattern, handler, mws, w, req)
	}
}
