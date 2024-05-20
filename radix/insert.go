package radix

import (
	"fmt"
)

func (node *Node) Insert(pattern string, value any) error {
	if node.code == 0 {
		node.code = '/'
	}
	if pattern == "/" {
		node.value = value
		return nil
	}
	input := normalizePath(pattern)
	if len(input) == 0 {
		return fmt.Errorf("pattern must be specified")
	}
	return createNodeRecursive(node, input, value, 0, 1)
}

func createNodeRecursive(root *Node, input []byte, value any, priority, segment uint8) error {
	if len(input) == 0 {
		return fmt.Errorf("input cannot be empty")
	}
	code := input[0]
	child := getChild(root, code)

	if code == '*' {
		if len(input) > 1 {
			return fmt.Errorf("wildcard is not supported in the middle of path")
		}
		if child != nil {
			return fmt.Errorf("attempt to redeclare wildcard value")
		}
		addChild(root, &Node{code: code, value: value, priority: priority + 1})
		return nil
	}

	if code == ':' {
		if len(input) == 1 {
			return fmt.Errorf("param value must be specified after ':'")
		}
		param, tail := extractParam(input)
		if len(tail) == 0 {
			if child != nil {
				if child.value != nil {
					return fmt.Errorf("try to redeclare node")
				}
				child.value = value
				child.priority = priority + 2
				child.segment = segment
				return nil
			}
			addChild(root, &Node{code: ':', param: string(param), segment: segment, priority: priority + 2, value: value})
			return nil
		}
		if child != nil {
			if child.param != string(param) {
				return fmt.Errorf("param name must be the same")
			}
		} else {
			child = addChild(root, &Node{code: ':', param: string(param), segment: segment})
		}
		return createNodeRecursive(child, tail, value, priority+1, segment)
	}

	// static node
	if code == '/' {
		segment++
	}

	if len(input) == 1 {
		if child != nil {
			if child.value != nil {
				return fmt.Errorf("try to redeclare node")
			}
			child.value = value
			return nil
		}
		addChild(root, &Node{code: code, priority: priority + 3, value: value})
		return nil
	}
	if child == nil {
		child = addChild(root, &Node{code: code})
	}
	return createNodeRecursive(child, input[1:], value, priority+1, segment)
}
