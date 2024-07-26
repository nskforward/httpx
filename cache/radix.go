package cache

import (
	"fmt"
	"strings"
	"sync"
)

type RadixNode struct {
	sync.RWMutex
	code     byte
	tag      *Tag
	children []*RadixNode
}

func (n *RadixNode) GetOrCreate(name []byte) *RadixNode {
	curr := name[0]
	tail := name[1:]

	if curr == '*' {
		panic(fmt.Errorf("wildcard symbol does not support in the cache tag name: %s", string(name)))
	}

	child := n.getChild(curr)
	if child == nil {
		child = n.createChild(curr)
	}
	if len(tail) == 0 {
		return child
	}
	return child.GetOrCreate(tail)
}

func (n *RadixNode) RangeChildren(f func(*RadixNode) bool) bool {
	for _, child := range n.children {
		if !f(child) {
			return false
		}
		if !child.RangeChildren(f) {
			return false
		}
	}
	return true
}

func (n *RadixNode) createChild(code byte) *RadixNode {
	node := &RadixNode{code: code}
	if n.children == nil {
		n.children = make([]*RadixNode, 0, 4)
	}
	n.children = append(n.children, node)
	return node
}

func (n *RadixNode) getChild(code byte) *RadixNode {
	for _, child := range n.children {
		if child.code == code {
			return child
		}
	}
	return nil
}

func (n *RadixNode) Dump(offset int) {
	for _, child := range n.children {
		fmt.Print(strings.Repeat("- ", offset))
		fmt.Println(string(child.code))
		child.Dump(offset + 1)
	}
}
