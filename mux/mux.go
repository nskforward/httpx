package mux

import (
	"fmt"
	"net/http"
)

type Multiplexer struct {
	root         *Node
	errorHandler func(w http.ResponseWriter, r *http.Request, code int)
}

func NewMultiplexer() *Multiplexer {
	return &Multiplexer{
		root: NewNode(nil, Token{Kind: Sep}),
		errorHandler: func(w http.ResponseWriter, r *http.Request, code int) {
			http.Error(w, http.StatusText(code), code)
		},
	}
}

func (m *Multiplexer) SetErrorHandler(h func(w http.ResponseWriter, r *http.Request, code int)) {
	m.errorHandler = h
}

func (m *Multiplexer) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	m.Handle(pattern, http.HandlerFunc(handler))
}

func (m *Multiplexer) Handle(pattern string, handler http.Handler) {
	method, url, err := SplitPattern(pattern)
	if err != nil {
		panic(err)
	}

	if len(url) == 0 || url[0] != '/' {
		panic("bad url format")
	}
	curr := m.root
	for segment := range Segments(url) {
		if segment == "/" {
			if curr == m.root {
				continue
			}
			curr = curr.GetChildByTokenOrCreate(Token{Kind: Sep})
			continue
		}
		if len(segment) > 2 && segment[0] == '{' && segment[len(segment)-1] == '}' {
			curr = curr.GetChildByTokenOrCreate(Token{Kind: Param, Param: segment[1 : len(segment)-1]})
			continue
		}
		for _, char := range segment {
			curr = curr.GetChildByTokenOrCreate(Token{Kind: Lit, Lit: char})
			continue
		}
	}
	curr.SetValue(method, handler)
}

func (m *Multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methodNum := MethodToUInt8(r.Method)
	child, exactly := m.root.GetLongest(r.URL.Path, r.SetPathValue)
	if child == nil {
		m.errorHandler(w, r, 404)
		return
	}
	if exactly {
		if child.value != nil {
			h := child.value.Get(methodNum)
			if h != nil {
				h.ServeHTTP(w, r)
				return
			}
			m.errorHandler(w, r, 405)
			return
		}
	}

	if child.wildcard != nil {
		h := child.wildcard.Get(methodNum)
		if h == nil {
			m.errorHandler(w, r, 405)
			return
		}
		h.ServeHTTP(w, r)
		return
	}
	m.errorHandler(w, r, 404)
}

func (m *Multiplexer) Dump() {
	fmt.Println("---------------------------")
	fmt.Println("# DUMP")
	fmt.Println("---------------------------")
	m.root.dump(0)
}
