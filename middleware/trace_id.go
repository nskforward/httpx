package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func TraceID(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := r.Header.Get(types.XTraceID)
		if id == "" {
			id = uuid.New().String()
			r.Header.Set(types.XTraceID, id)
		}

		ww := response.NewWrapper(w)
		ww.BeforeBody = func() {
			if ww.Header().Get(types.XTraceID) == "" {
				ww.Header().Set(types.XTraceID, id)
			}
		}

		return next(ww, r)
	}
}
