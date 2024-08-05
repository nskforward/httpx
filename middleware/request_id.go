package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nskforward/httpx/types"
)

func RequestID(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := r.Header.Get(types.XRequestId)
		needsToSend := false
		if id == "" {
			needsToSend = true
			id = uuid.New().String()
			r.Header.Set(types.XRequestId, id)
		}

		ww := types.NewResponseWrapper(w)
		ww.BeforeBody = func() {
			if ww.Header().Get(types.XRequestId) == "" && needsToSend {
				ww.Header().Set(types.XRequestId, id)
			}
		}

		return next(ww, r)
	}
}
