package httpx

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/nskforward/httpx/logging"
)

type LoggingFunc func(w *logging.Writer, r *http.Request)

func DefaultLoggingFunc(w *logging.Writer, r *http.Request) {
	slog.Info("incoming request",
		"status", w.StatusCode(),
		"method", r.Method,
		"path", r.URL.Path,
		"proto", r.Proto,
		"trace", w.Header().Get(TraceIDHeader),
		"client", UserIP(r),
		"spend", w.Duration().String(),
		"bytes", strconv.Itoa(w.Size()))
}
