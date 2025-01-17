package httpx

import (
	"fmt"
	"strings"
)

/*
Patterns can match the method, host and path of a request. Some examples:

	"/index.html" matches the path "/index.html" for any host and method.
	"GET /static/" matches a GET request whose path begins with "/static/".
	"example.com/" matches any request to the host "example.com".
	"example.com/{$}" matches requests with host "example.com" and path "/".
	"/b/{bucket}/o/{objectname...}" matches paths whose first segment is "b" and whose third segment is "o".
*/
func DeclareHandler(r *Router, method, pattern string, handler Handler, middlewares ...Middleware) {
	mws := append(r.mws, middlewares...)

	if method == "" {
		if strings.HasPrefix(pattern, "/") {
			r.serverMux.HandleFunc(pattern, toStdHandler(r, pattern, handler, mws))
			return
		}
		panic(fmt.Errorf("router: handler declaration: invalid pattern for method ANY: %s", pattern))
	}

	if strings.HasPrefix(pattern, fmt.Sprintf("%s ", method)) {
		r.serverMux.HandleFunc(pattern, toStdHandler(r, pattern, handler, mws))
		return
	}

	if strings.HasPrefix(pattern, "/") {
		pattern = fmt.Sprintf("%s %s", method, pattern)
		r.serverMux.HandleFunc(pattern, toStdHandler(r, pattern, handler, mws))
		return
	}

	panic(fmt.Errorf("router: handler declaration: invalid pattern for method %s: %s", method, pattern))
}
