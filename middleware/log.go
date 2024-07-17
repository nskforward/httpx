package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func Log(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ww := response.NewWrapper(w)

		t1 := time.Now()
		err := catchError(next)(ww, r)
		spent := time.Since(t1)

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		slog.Info("access",
			"status", ww.Status(),
			"method", r.Method,
			"path", r.URL.Path,
			"proto", r.Proto,
			"request-id", r.Header.Get(types.XRequestId),
			"client", ip,
			"spend", spent.String(),
			"bytes", ww.Size(),
		)

		return err
	}
}

func catchError(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := next(w, r)
		if err == nil {
			return nil
		}

		resp, ok := err.(response.Error)
		if ok {
			slog.Error(strconv.Itoa(resp.Status), "error", resp.Text, "trace-id", r.Header.Get(types.XRequestId))
			if resp.Text == "" || resp.Status == 500 {
				http.Error(w, http.StatusText(resp.Status), resp.Status)
			} else {
				http.Error(w, resp.Text, resp.Status)
			}
		} else {
			slog.Error("400", "error", err.Error(), "request-id", r.Header.Get(types.XRequestId))
			http.Error(w, err.Error(), 400)
		}

		return nil
	}
}
