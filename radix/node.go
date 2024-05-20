package radix

import (
	"bytes"
	"unsafe"
)

type Node struct {
	parent   *Node
	priority uint8
	segment  uint8
	code     byte
	param    string
	value    any
	children []*Node
}

func (n *Node) Value() any {
	return n.value
}

func getChild(parent *Node, code byte) *Node {
	for _, child := range parent.children {
		if child.code == code {
			return child
		}
	}
	return nil
}

func addChild(parent *Node, child *Node) *Node {
	if parent.children == nil {
		parent.children = make([]*Node, 0, 4)
	}
	child.parent = parent
	parent.children = append(parent.children, child)
	return child
}

func normalizePath(s string) []byte {
	return bytes.Trim(s2b(s), "/")
}

func s2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func b2s(b []byte) string {
	return unsafe.String(&b[0], len(b))
}

func extractParam(input []byte) (param []byte, tail []byte) {
	if len(input) == 0 {
		return nil, nil
	}
	if input[0] == ':' {
		input = input[1:]
	}
	pos := bytes.Index(input, []byte("/"))
	if pos < 0 {
		return input, nil
	}
	return input[:pos], input[pos:]
}
