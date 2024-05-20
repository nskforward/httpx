package radix

func (node *Node) Search(path string) *Node {
	if path == "/" {
		if node.value != nil {
			return node
		}
		return nil
	}
	input := normalizePath(path)
	if len(input) == 0 {
		return nil
	}
	result := findBestChildren(node, input)
	if result != nil && result.value != nil {
		return result
	}
	return nil
}

func findBestChildren(root *Node, input []byte) *Node {
	if len(input) == 0 {
		return nil
	}
	bestNode := chooseBestNode(searchStaticNode(root, input), searchParamNode(root, input))
	if bestNode != nil {
		return bestNode
	}
	return searchWildcardNode(root)
}

func chooseBestNode(n1, n2 *Node) *Node {
	if n1 == nil && n2 == nil {
		return nil
	}
	if n1 == nil && n2 != nil {
		return n2
	}
	if n1 != nil && n2 == nil {
		return n1
	}
	if n1.priority > n2.priority {
		return n1
	}
	return n2
}

func searchStaticNode(parent *Node, input []byte) *Node {
	child := getChild(parent, input[0])
	if child != nil {
		if len(input) == 1 {
			return child
		}
		return findBestChildren(child, input[1:])
	}
	return nil
}

func searchParamNode(parent *Node, input []byte) *Node {
	child := getChild(parent, ':')
	if child == nil {
		return nil
	}
	_, tail := extractParam(input)
	if len(tail) == 0 && child.value != nil {
		return child
	}
	return findBestChildren(child, tail)
}

func searchWildcardNode(parent *Node) *Node {
	child := getChild(parent, '*')
	if child != nil {
		return child
	}
	return nil
}
