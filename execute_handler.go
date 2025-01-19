package httpx

import (
	"fmt"
	"net/http"
	"slices"
)

func executeHandler(r *Router, pattern string, h Handler, middlewares []Middleware, w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(r.logger, pattern, w, req)

	if r.beforeRequestLog != nil {
		r.beforeRequestLog(ctx)
	}

	finalHandler := h
	for _, mw := range slices.Backward(middlewares) {
		finalHandler = mw(finalHandler)
	}

	err := finalHandler(ctx)
	if err != nil {
		ctx.Logger().Error("unexpected", "error", err)
		if !ctx.HeadersSent() {
			ctx.RespondText(http.StatusInternalServerError, fmt.Sprintf("internal server error: trace id: %s", ctx.TraceID()))
		}
	}

	if r.afterResponseLog != nil {
		r.afterResponseLog(ctx)
	}
}
