package httpx

import (
	"log/slog"
	"net/http"

	"github.com/nskforward/httpx/types"
)

func defaultErrorFunc(w http.ResponseWriter, r *http.Request, status int, msg string) {
	slog.Error(msg, "trace-id", r.Header.Get(types.XRequestId))

	if status/100 == 5 {
		http.Error(w, http.StatusText(status), status)
		return
	}

	http.Error(w, msg, status)
}
