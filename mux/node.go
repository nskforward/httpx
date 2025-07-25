package mux

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type Node struct {
	tok      Token
	parent   *Node
	children []*Node

	wildcard *Leaf
	value    *Leaf
}

func NewNode(parent *Node, tok Token) *Node {
	return &Node{
		parent:   parent,
		children: make([]*Node, 0, 8),
		tok:      tok,
	}
}

func (n *Node) GetLongest(input string, setParam func(name, value string)) (*Node, bool) {
	var best *Node
	curr := n
	for segment := range Segments(input) {
		if segment == "/" && curr == n {
			if n.wildcard != nil {
				best = n
			}
			continue
		}
		for _, char := range segment {
			child := curr.GetChildByRune(char)
			if child == nil {
				return best, false
			}
			if child.wildcard != nil {
				best = child
			}
			if child.tok.Kind == Param && child.tok.Param != "$" {
				setParam(child.tok.Param, segment)
				curr = child
				break
			}
			curr = child
		}
	}
	return curr, true
}

func (n *Node) SetValue(method Method, handler http.Handler) error {
	if n.tok.Kind == Lit {
		return n.setValueExactly(true, method, handler)
	}
	if n.tok.Kind == Sep {
		return n.setValueExactly(false, method, handler)
	}
	if n.tok.Kind == Param {
		if n.tok.Param == "$" && n.parent != nil && n.parent.tok.Kind == Sep {
			return n.parent.setValueExactly(true, method, handler)
		}
		return n.setValueExactly(true, method, handler)
	}
	return fmt.Errorf("cannot set handler on node kind %s", n.tok.Kind.String())
}

func (n *Node) setValueExactly(exactly bool, method Method, handler http.Handler) error {
	if exactly {
		if n.value == nil {
			n.value = newLeaf()
		}
		return n.value.Set(method, handler)
	}
	if n.wildcard == nil {
		n.wildcard = newLeaf()
	}
	return n.wildcard.Set(method, handler)
}

func (n *Node) GetChildByTokenOrCreate(token Token) *Node {
	child := n.GetChildByToken(token)
	if child == nil {
		child = NewNode(n, token)
		n.children = append(n.children, child)
	}
	return child
}

func (n *Node) GetChildByToken(token Token) *Node {
	for _, child := range n.children {
		if child.tok.Kind != token.Kind {
			continue
		}
		if token.Kind == Lit && token.Lit == child.tok.Lit {
			return child
		}
		if token.Kind == Param {
			return child
		}
		if token.Kind == Sep {
			return child
		}
	}
	return nil
}

func (n *Node) GetChildByRune(char rune) *Node {
	var best *Node
	for _, child := range n.children {
		if char == '/' && child.tok.Kind == Sep {
			return child
		}
		if child.tok.Kind == Lit && child.tok.Lit == char {
			return child
		}
		if child.tok.Kind == Param {
			best = child
		}
	}
	return best
}

func (n *Node) String() string {
	return n.tok.String()
}

func (n *Node) dump(offset int) {
	var buf bytes.Buffer
	if n.value != nil {
		buf.WriteString("exact=[")
		buf.WriteString(n.value.String())
		buf.WriteString("] ")
	}
	if n.wildcard != nil {
		buf.WriteString("wildcard=[")
		buf.WriteString(n.wildcard.String())
		buf.WriteString("] ")
	}

	fmt.Printf("%s%s %s\n", strings.Repeat("  ", offset), n.tok, buf.String())
	for _, child := range n.children {
		child.dump(offset + 1)
	}
}
