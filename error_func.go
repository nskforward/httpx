package httpx

import (
	"log/slog"
	"net/http"
)

func DefaultErrorFunc(w http.ResponseWriter, r *http.Request, status int, msg string) {

	slog.Error(msg, "trace-id", TraceID(r.Context()))

	if status/100 == 5 {
		http.Error(w, http.StatusText(status), status)
		return
	}

	http.Error(w, msg, status)
}
