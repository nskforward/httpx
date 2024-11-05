package httpx

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/nskforward/httpx/middleware"
	"github.com/nskforward/httpx/types"
)

func DefaultLoggerFunc(w *types.ResponseWrapper, r *http.Request) {

	slog.Info(r.Method,
		"status", w.Status(),
		"path", r.URL.Path,
		"proto", r.Proto,
		"trace-id", middleware.GetTraceID(r.Context()),
		"peer", r.RemoteAddr,
		"spent", time.Since(w.StartTime()).String(),
		"bytes", w.Size(),
	)
}
