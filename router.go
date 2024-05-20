package httpx

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/nskforward/httpx/logging"
)

type Router struct {
	mu           ServeMux
	loggingFunc  LoggingFunc
	defaultGroup Group
}

func NewRouter() *Router {
	r := &Router{
		loggingFunc: DefaultLoggingFunc,
	}
	r.defaultGroup = Group{router: r, middlewares: make([]Middleware, 0, 16)}
	return r
}

func (ro *Router) NewGroup(middleware ...Middleware) Group {
	return Group{
		router:      ro,
		middlewares: append(ro.defaultGroup.middlewares, middleware...),
	}
}

func (ro *Router) UseLogging(loggingFunc LoggingFunc) {
	ro.loggingFunc = loggingFunc
}

func (ro *Router) Use(middlewares ...Middleware) {
	ro.defaultGroup.middlewares = append(ro.defaultGroup.middlewares, middlewares...)
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ww := logging.Get()
	ww.Reset(w)

	defer logging.Put(ww)
	defer ww.Close()

	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			buf = append(buf, '.', '.', '.')
			msg := fmt.Sprintf("%v", err)

			InternalServerError(fmt.Errorf("panic: %s", msg)).ServeHTTP(w, r)

			slog.Error("panic", "error", msg, "stacktrace", string(buf))
		}
	}()

	ro.mu.ServeHTTP(ww, r)

	ro.loggingFunc(ww, r)
}

func (ro *Router) ANY(pattern string, h Handler, middlewares ...Middleware) {
	ro.defaultGroup.ANY(pattern, h, middlewares...)
}

func (ro *Router) GET(pattern string, h Handler, middlewares ...Middleware) {
	ro.defaultGroup.GET(pattern, h, middlewares...)
}

func (ro *Router) POST(pattern string, h Handler, middlewares ...Middleware) {
	ro.defaultGroup.POST(pattern, h, middlewares...)
}

func (ro *Router) DELETE(pattern string, h Handler, middlewares ...Middleware) {
	ro.defaultGroup.DELETE(pattern, h, middlewares...)
}

func (ro *Router) PUT(pattern string, h Handler, middlewares ...Middleware) {
	ro.defaultGroup.PUT(pattern, h, middlewares...)
}

func (ro *Router) PATCH(pattern string, h Handler, middlewares ...Middleware) {
	ro.defaultGroup.PATCH(pattern, h, middlewares...)
}

func (ro *Router) registryRoute(method string, route Route) {
	h := route.handler
	for i := len(route.middlewares) - 1; i >= 0; i-- {
		h = route.middlewares[i](h)
	}
	ro.mu.Route(method, route.pattern, ro.Catch(h))
}
