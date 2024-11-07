package httpx

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/nskforward/httpx/types"
)

var contextTraceID types.ContextParam = "middleware.trace.id"

func traceIDSetter(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {

		id := r.Header.Get(types.XTraceID)
		if id == "" {
			id = uuid.New().String()
		}

		w.Header().Set(types.XTraceID, id)
		r = types.SetParam(r, contextTraceID, id)

		return next(w, r)
	}
}

func TraceID(ctx context.Context) string {
	id := ctx.Value(contextTraceID)
	if id == nil {
		return ""
	}
	return id.(string)
}
