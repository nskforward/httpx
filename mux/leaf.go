package mux

import (
	"fmt"
	"net/http"
	"strings"
)

type Leaf struct {
	methods map[Method]http.Handler
}

func newLeaf() *Leaf {
	return &Leaf{
		methods: make(map[Method]http.Handler),
	}
}

func (l *Leaf) Set(method Method, handler http.Handler) error {
	_, ok := l.methods[method]
	if ok {
		return fmt.Errorf("handler already defined for method %s", MethodToStr(method))
	}
	l.methods[method] = handler
	return nil
}

func (l *Leaf) Get(method Method) http.Handler {
	h, ok := l.methods[method]
	if !ok && method != ANY {
		return l.methods[ANY]
	}
	return h
}

func (l *Leaf) String() string {
	methods := make([]string, 0, len(l.methods))
	for method := range l.methods {
		methods = append(methods, MethodToStr(method))
	}
	return strings.Join(methods, ", ")
}
