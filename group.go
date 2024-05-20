package httpx

type Group struct {
	router      *Router
	middlewares []Middleware
}

func (g Group) NewGroup(middlewares ...Middleware) Group {
	return Group{
		router:      g.router,
		middlewares: append(g.middlewares, middlewares...),
	}
}

func (g Group) ANY(pattern string, handler Handler, middlewares ...Middleware) {
	g.router.registryRoute("ANY", Route{
		pattern:     pattern,
		handler:     handler,
		middlewares: append(g.middlewares, middlewares...),
	})
}

func (g Group) GET(pattern string, handler Handler, middlewares ...Middleware) {
	g.router.registryRoute("GET", Route{
		pattern:     pattern,
		handler:     GET(handler),
		middlewares: append(g.middlewares, middlewares...),
	})
}

func (g Group) POST(pattern string, handler Handler, middlewares ...Middleware) {
	g.router.registryRoute("POST", Route{
		pattern:     pattern,
		handler:     POST(handler),
		middlewares: append(g.middlewares, middlewares...),
	})
}

func (g Group) PUT(pattern string, handler Handler, middlewares ...Middleware) {
	g.router.registryRoute("PUT", Route{
		pattern:     pattern,
		handler:     PUT(handler),
		middlewares: append(g.middlewares, middlewares...),
	})
}

func (g Group) PATCH(pattern string, handler Handler, middlewares ...Middleware) {
	g.router.registryRoute("PATCH", Route{
		pattern:     pattern,
		handler:     PATCH(handler),
		middlewares: append(g.middlewares, middlewares...),
	})
}

func (g Group) DELETE(pattern string, handler Handler, middlewares ...Middleware) {
	g.router.registryRoute("DELETE", Route{
		pattern:     pattern,
		handler:     DELETE(handler),
		middlewares: append(g.middlewares, middlewares...),
	})
}
