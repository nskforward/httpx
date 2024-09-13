package httpx

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/nskforward/httpx/types"
)

func defaultLoggerFunc(w *types.ResponseWrapper, r *http.Request, err error) {
	if err != nil {
		slog.Error(strconv.Itoa(w.Status()), "error", err, "trace-id", r.Header.Get(types.XRequestId))
		handleError(w, err)
	}

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
