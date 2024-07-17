package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func Logging(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		ww := response.NewWrapper(w)

		t1 := time.Now()
		err = next(ww, r)
		spent := time.Since(t1)

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		slog.Info("incoming request",
			"status", ww.Status(),
			"method", r.Method,
			"path", r.URL.Path,
			"proto", r.Proto,
			"trace-id", r.Header.Get(types.XTraceID),
			"client", ip,
			"spend", spent.String(),
			"bytes", ww.Size(),
		)

		return
	}
}
