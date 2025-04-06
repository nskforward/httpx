package httpx

type Group struct {
	router      *Router
	middlewares []Handler
	pattern     string
}

func NewGroup(router *Router, pattern string, middlewares []Handler) *Group {
	return &Group{
		router:      router,
		pattern:     pattern,
		middlewares: middlewares,
	}
}

func (g *Group) Use(middleware ...Handler) {
	g.middlewares = append(g.middlewares, middleware...)
}

func (g *Group) Group(pattern string, middlewares ...Handler) *Group {
	return NewGroup(g.router, g.joinPattern(pattern), append(g.middlewares, middlewares...))
}

func (g *Group) joinPattern(pattern string) string {
	if len(g.pattern) == 0 || len(pattern) == 0 {
		return g.pattern + pattern
	}
	if g.pattern[len(g.pattern)-1] == '/' && pattern[0] == '/' {
		return g.pattern + pattern[1:]
	}
	return g.pattern + pattern
}

func (g *Group) DELETE(pattern string, handler Handler, middleware ...Handler) {
	g.CustomMethod("DELETE", pattern, handler, middleware...)
}

func (g *Group) PATCH(pattern string, handler Handler, middleware ...Handler) {
	g.CustomMethod("PATCH", pattern, handler, middleware...)
}

func (g *Group) PUT(pattern string, handler Handler, middleware ...Handler) {
	g.CustomMethod("PUT", pattern, handler, middleware...)
}

func (g *Group) POST(pattern string, handler Handler, middleware ...Handler) {
	g.CustomMethod("POST", pattern, handler, middleware...)
}

func (g *Group) GET(pattern string, handler Handler, middleware ...Handler) {
	g.CustomMethod("GET", pattern, handler, middleware...)
}

func (g *Group) ANY(pattern string, handler Handler, middleware ...Handler) {
	g.CustomMethod("", pattern, handler, middleware...)
}

func (g *Group) CustomMethod(method, pattern string, handler Handler, middleware ...Handler) {
	g.router.CustomMethod(method, g.joinPattern(pattern), handler, append(g.middlewares, middleware...)...)
}
