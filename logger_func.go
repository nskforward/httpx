package httpx

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/nskforward/httpx/types"
)

func defaultLoggerFunc(w *types.ResponseWrapper, r *http.Request) {

	slog.Info(r.Method,
		"status", w.Status(),
		"path", r.URL.Path,
		"proto", r.Proto,
		"request-id", r.Header.Get(types.XRequestId),
		"peer", r.RemoteAddr,
		"spent", time.Since(w.StartTime()).String(),
		"bytes", w.Size(),
	)
}
