package radix

import (
	"fmt"
)

func DumpTree(node *Node, level int) {
	dumpPrintPadding(level)
	fmt.Println(dumpParam(node), dumpValue(node))
	for _, child := range node.children {
		DumpTree(child, level+1)
	}
}

func dumpPrintPadding(level int) {
	for i := 0; i < level; i++ {
		if i%2 == 0 {
			fmt.Print(" .")
		} else {
			fmt.Print("  ")
		}
	}
}

func dumpParam(node *Node) string {
	if node.code != ':' {
		return fmt.Sprintf("[%c]", node.code)
	}
	return fmt.Sprintf("[:%s/%d]", node.param, node.segment)
}

func dumpValue(node *Node) string {
	if node.value == nil {
		return ""
	}
	return fmt.Sprintf("-%d- (%v)", node.priority, node.value)
}

/*
func DumpChain(last *Node) string {
	var buf bytes.Buffer
	curr := last
	for curr != nil {
		if curr.code == ':' {
			buf.WriteString(reverseString(curr.param))
		}
		buf.WriteByte(curr.code)
		curr = curr.parent
	}
	b := buf.Bytes()
	slices.Reverse(b)
	return string(b)
}
*/
