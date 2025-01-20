package httpx

import (
	"log/slog"
	"time"
)

func DefaultLoggingFunc(ctx *Context) {

	requestHeaders := make([]any, 0, 32)
	for name, values := range ctx.Request().Header {
		for _, value := range values {
			requestHeaders = append(requestHeaders, name)
			requestHeaders = append(requestHeaders, value)
		}
	}
	responseHeaders := make([]any, 0, 32)
	for name, values := range ctx.ResponseWriter().Header() {
		for _, value := range values {
			responseHeaders = append(responseHeaders, name)
			responseHeaders = append(responseHeaders, value)
		}
	}

	ctx.Logger().Info(
		"http_call",
		slog.Group("request",
			slog.String("method", ctx.Method()),
			slog.String("path", ctx.Request().RequestURI),
			slog.String("client_ip", ctx.UserIP()),
			slog.Group("headers", requestHeaders...),
		),
		slog.Group("response",
			slog.Int("status", ctx.StatusCode()),
			slog.Int64("duration_ms", time.Since(ctx.StartTime()).Milliseconds()),
			slog.Int64("size", ctx.ResponseSize()),
			slog.Group("headers", responseHeaders...),
		),
	)
}
