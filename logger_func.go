package httpx

import (
	"log/slog"
	"net/http"

	"github.com/nskforward/httpx/types"
)

func DefaultLogger(w *types.ResponseWrapper, r *http.Request) {

	slog.Info("http request",
		"method", r.Method,
		"status", w.Status(),
		"path", r.URL.Path,
		"proto", r.Proto,
		"trace-id", TraceID(r.Context()),
		"peer", r.RemoteAddr,
		"spent", w.TimeTaken().String(),
		"bytes", w.Size(),
	)
}
