package httpx

import (
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/google/uuid"
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
		logger := ro.logger.With("id", GetRequestID(r))
		resp := newResponse(w, logger)
		err := final(resp, r)
		if err != nil {
			ro.errorHandler(resp, r, err)
		}
	}
}

func defaultErrorHandler(w *Response, r *http.Request, err error) {
	fmt.Fprintf(os.Stderr, "ERROR: httpx: unhandler error during the route '%s %s': %v\n", r.Method, r.URL.Path, err)
	if !w.HeadersSent() {
		w.SendShortError(http.StatusInternalServerError)
	}
}

func GetRequestID(r *http.Request) string {
	requestID := r.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = uuid.NewString()
		r.Header.Set("X-Request-Id", requestID)
	}
	return requestID
}
