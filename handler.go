package httpx

import (
	"net/http"
	"slices"
)

func castHandler(ro *Router, h HandlerFunc, mws []Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		final := h
		for _, mw := range slices.Backward(mws) {
			final = mw(final)
		}
		for _, mw := range slices.Backward(ro.middlewares) {
			final = mw(final)
		}
		resp := newResponse(w)
		err := final(resp, r)
		if err != nil {
			ro.handlerError(resp, r, err)
		}
	}
}
