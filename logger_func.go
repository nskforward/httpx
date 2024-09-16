package httpx

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/nskforward/httpx/types"
)

func defaultLoggerFunc(w *types.ResponseWrapper, r *http.Request) {

	slog.Info("access",
		"status", w.Status(),
		"method", r.Method,
		"path", r.URL.Path,
		"proto", r.Proto,
		"request-id", r.Header.Get(types.XRequestId),
		"peer", r.RemoteAddr,
		"spend", time.Since(w.StartTime()).String(),
		"bytes", w.Size(),
	)
}
