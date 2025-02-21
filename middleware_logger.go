package httpx

import (
	"log/slog"
	"time"
)

func Logger(logHeaders bool) Handler {
	return func(ctx *Ctx) error {

		t1 := time.Now()
		err := ctx.Next()
		t2 := time.Since(t1)

		var requestHeaders, responseHeaders []any

		if logHeaders {
			requestHeaders = make([]any, 0, 32)
			for name, values := range ctx.Request().Header {
				for _, value := range values {
					requestHeaders = append(requestHeaders, name)
					requestHeaders = append(requestHeaders, value)
				}
			}
			responseHeaders = make([]any, 0, 32)
			for name, values := range ctx.w.Header() {
				for _, value := range values {
					responseHeaders = append(responseHeaders, name)
					responseHeaders = append(responseHeaders, value)
				}
			}
		}

		ctx.Logger().Info(
			"http_call",
			slog.Group("request",
				slog.String("method", ctx.Request().Method),
				slog.String("path", ctx.Request().RequestURI),
				slog.Group("user",
					slog.String("ip", ctx.clientAddr),
					slog.String("agent", ctx.Request().UserAgent()),
				),
				slog.Group("headers", requestHeaders...),
			),
			slog.Group("response",
				slog.Int("status", ctx.w.Status()),
				slog.Int64("duration_ms", t2.Milliseconds()),
				slog.Int64("size", ctx.w.Size()),
				slog.Group("headers", responseHeaders...),
			),
		)

		return err
	}

}
