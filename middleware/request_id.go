package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nskforward/httpx/types"
)

func RequestID(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := r.Header.Get(types.XRequestId)
		if id == "" {
			id = r.Header.Get("Cf-Ray")
			if id != "" {
				r.Header.Set(types.XRequestId, id)
			}
		}

		if id == "" {
			id = uuid.New().String()
			r.Header.Set(types.XRequestId, id)
			w.Header().Set(types.XRequestId, id)
		}

		return next(w, r)
	}
}
