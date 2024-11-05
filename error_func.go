package httpx

import (
	"log/slog"
	"net/http"

	"github.com/nskforward/httpx/middleware"
)

func DefaultErrorFunc(w http.ResponseWriter, r *http.Request, status int, msg string) {
	slog.Error(msg, "trace-id", middleware.GetTraceID(r.Context()))

	if status/100 == 5 {
		http.Error(w, http.StatusText(status), status)
		return
	}

	http.Error(w, msg, status)
}
