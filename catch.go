package httpx

import (
	"log/slog"
	"net/http"
)

func (ro *Router) Catch(next Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)

		if err == nil {
			return
		}

		slog.Debug("handling finished with an error", "trace", GetTraceID(r), "error", err)

		entity, ok := err.(ResponseError)
		if !ok {
			entity = BadRequest(err)
		}

		entity.ServeHTTP(w, r)
	}
}
