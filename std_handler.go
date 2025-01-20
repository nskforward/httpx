package httpx

import (
	"fmt"
	"net/http"
	"slices"
)

func wrapHandler(router *Router, handler Handler, mws []Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		executeFinalHandler(router, handler, mws, w, r)
	}
}

func executeFinalHandler(router *Router, handler Handler, middlewares []Middleware, w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(router.logger, w, r)

	if router.beforeRequestLog != nil {
		router.beforeRequestLog(ctx)
	}

	finamHandler := handler
	for _, mw := range slices.Backward(middlewares) {
		finamHandler = mw(finamHandler)
	}

	err := finamHandler(ctx)
	if err != nil {
		ctx.Logger().Error("unexpected", "error", err)
		if !ctx.HeadersSent() {
			ctx.RespondText(http.StatusInternalServerError, fmt.Sprintf("internal server error: trace id: %s", ctx.TraceID()))
		}
	}

	if router.afterResponseLog != nil {
		router.afterResponseLog(ctx)
	}
}
