package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/nskforward/httpx"
)

func Logger(logHeaders bool) httpx.Handler {
	return func(req *http.Request, resp *httpx.Response) error {

		t1 := time.Now()
		err := resp.Next()
		t2 := time.Since(t1)

		if err != nil {
			resp.Logger().Error(err.Error())
			resp.InternalServerError()
		}

		var requestHeaders, responseHeaders []any

		if logHeaders {
			requestHeaders = make([]any, 0, 32)
			for name, values := range req.Header {
				for _, value := range values {
					requestHeaders = append(requestHeaders, name)
					requestHeaders = append(requestHeaders, value)
				}
			}
			responseHeaders = make([]any, 0, 32)
			for name, values := range resp.ResponseWriter().Header() {
				for _, value := range values {
					responseHeaders = append(responseHeaders, name)
					responseHeaders = append(responseHeaders, value)
				}
			}
		}

		resp.Logger().Info(
			"http_call",
			slog.Group("request",
				slog.String("method", req.Method),
				slog.String("path", req.RequestURI),
				slog.Group("user",
					slog.String("ip", req.RemoteAddr),
					slog.String("agent", req.UserAgent()),
				),
				slog.Group("headers", requestHeaders...),
			),
			slog.Group("response",
				slog.Int("status", resp.StatusCode()),
				slog.Int64("duration_ms", t2.Milliseconds()),
				slog.Int64("size", resp.BodySize()),
				slog.Group("headers", responseHeaders...),
			),
		)

		return nil
	}

}
