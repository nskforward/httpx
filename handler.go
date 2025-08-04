package httpx

import (
	"log/slog"
	"net/http"
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
		var logger *slog.Logger
		if ro.logger != nil {
			logger = ro.logger.With("id", GetRequestID(r))
		}

		resp := newResponse(w, logger)
		err := final(resp, r)
		if err != nil {
			ro.handlerError(resp, r, err)
		}
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
