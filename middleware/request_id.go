package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func RequestID(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := r.Header.Get(types.XRequestId)
		if id == "" {
			id = uuid.New().String()
			r.Header.Set(types.XRequestId, id)
		}

		ww := response.NewWrapper(w)
		ww.BeforeBody = func() {
			if ww.Header().Get(types.XRequestId) == "" {
				ww.Header().Set(types.XRequestId, id)
			}
		}

		return next(ww, r)
	}
}
