package httpx

type Group struct {
	r           *router
	path        string
	middlewares []Handler
}

func (g *Group) Use(middlewares ...Handler) {
	g.middlewares = g.joinMiddleware(middlewares)
}

func (g *Group) Group(pattern string, middlewares ...Handler) Group {
	return Group{
		r:           g.r,
		path:        g.joinPath(pattern),
		middlewares: g.joinMiddleware(middlewares),
	}
}

func (g *Group) ANY(pattern string, handler Handler, middlewares ...Handler) {
	g.Custom("", pattern, handler, middlewares...)
}

func (g *Group) GET(pattern string, handler Handler, middlewares ...Handler) {
	g.Custom("GET", pattern, handler, middlewares...)
}

func (g *Group) POST(pattern string, handler Handler, middlewares ...Handler) {
	g.Custom("POST", pattern, handler, middlewares...)
}

func (g *Group) PUT(pattern string, handler Handler, middlewares ...Handler) {
	g.Custom("PUT", pattern, handler, middlewares...)
}

func (g *Group) DELETE(pattern string, handler Handler, middlewares ...Handler) {
	g.Custom("DELETE", pattern, handler, middlewares...)
}

func (g *Group) PATCH(pattern string, handler Handler, middlewares ...Handler) {
	g.Custom("PATCH", pattern, handler, middlewares...)
}

func (g *Group) OPTIONS(pattern string, handler Handler, middlewares ...Handler) {
	g.Custom("OPTIONS", pattern, handler, middlewares...)
}

func (g *Group) Custom(method, pattern string, handler Handler, middlewares ...Handler) {
	g.r.Route(method, g.joinPath(pattern), handler, g.joinMiddleware(middlewares))
}

func (g *Group) joinMiddleware(mw []Handler) []Handler {
	return append(g.middlewares, mw...)
}

func (g *Group) joinPath(path string) string {
	if len(g.path) == 0 || len(path) == 0 {
		return g.path + path
	}
	if g.path[len(g.path)-1] == '/' && path[0] == '/' {
		return g.path + path[1:]
	}
	return g.path + path
}
