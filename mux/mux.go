package mux

import (
	"fmt"
	"net/http"
)

type Multiplexer struct {
	root *Node
}

func NewMultiplexer() *Multiplexer {
	return &Multiplexer{
		root: NewNode(nil, Token{Kind: Sep}),
	}
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

func (m *Multiplexer) Search(w http.ResponseWriter, r *http.Request) (http.Handler, int) {
	return m.search(r)
}

func (m *Multiplexer) search(r *http.Request) (http.Handler, int) {
	methodNum := MethodToUInt8(r.Method)
	child, exactly := m.root.GetLongest(r.URL.Path, r.SetPathValue)
	if child == nil {
		return nil, 404
	}
	if exactly {
		if child.value != nil {
			h := child.value.Get(methodNum)
			if h != nil {
				return h, 0
			}
			return nil, 405
		}
	}
	if child.wildcard != nil {
		h := child.wildcard.Get(methodNum)
		if h == nil {
			return nil, 405
		}
		return h, 0
	}
	return nil, 404
}

func (m *Multiplexer) Dump() {
	fmt.Println("---------------------------")
	fmt.Println("# DUMP")
	fmt.Println("---------------------------")
	m.root.dump(0)
}
