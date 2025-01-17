package httpx

import "strings"

type Group struct {
	r             *Router
	mws           []Middleware
	patternPrefix string
}

func NewGroup(r *Router, patternPrefix string, middleware ...Middleware) *Group {
	return &Group{
		r:             r,
		mws:           middleware,
		patternPrefix: patternPrefix,
	}
}

func (g *Group) Use(middleware ...Middleware) {
	g.mws = append(g.mws, middleware...)
}

func (g *Group) Group(patternPrefix string) *Group {
	return NewGroup(g.r, strings.Join([]string{g.patternPrefix, patternPrefix}, ""), g.r.mws...)
}

func (g *Group) DeclareHandler(method, pattern string, handler Handler, middlewares ...Middleware) {
	DeclareHandler(g.r, "", strings.Join([]string{g.patternPrefix, pattern}, ""), handler, append(g.mws, middlewares...)...)
}

func (g *Group) ANY(pattern string, handler Handler, middlewares ...Middleware) {
	g.DeclareHandler("", pattern, handler, middlewares...)
}

func (g *Group) GET(pattern string, handler Handler, middlewares ...Middleware) {
	g.DeclareHandler("GET", pattern, handler, middlewares...)
}

func (g *Group) POST(pattern string, handler Handler, middlewares ...Middleware) {
	g.DeclareHandler("POST", pattern, handler, middlewares...)
}

func (g *Group) PUT(pattern string, handler Handler, middlewares ...Middleware) {
	g.DeclareHandler("PUT", pattern, handler, middlewares...)
}

func (g *Group) DELETE(pattern string, handler Handler, middlewares ...Middleware) {
	g.DeclareHandler("DELETE", pattern, handler, middlewares...)
}

func (g *Group) PATCH(pattern string, handler Handler, middlewares ...Middleware) {
	g.DeclareHandler("PATCH", pattern, handler, middlewares...)
}

func (g *Group) OPTIONS(pattern string, handler Handler, middlewares ...Middleware) {
	g.DeclareHandler("OPTIONS", pattern, handler, middlewares...)
}

func (g *Group) HEAD(pattern string, handler Handler, middlewares ...Middleware) {
	g.DeclareHandler("HEAD", pattern, handler, middlewares...)
}
