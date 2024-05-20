package radix

import "bytes"

func (node *Node) GetParam(path string, name string) string {
	input := normalizePath(path)
	curr := node
	for curr != nil {
		if curr.param == name {
			return getSegment(input, curr.segment)
		}
		curr = curr.parent
	}
	return ""
}

func getSegment(input []byte, segment uint8) string {
	if len(input) == 0 || input[0] == '/' {
		return ""
	}
	if segment < 1 {
		return ""
	}
	var pos1 int

	if segment == 1 {
		pos1 = bytes.IndexByte(input, '/')
		if pos1 < 0 {
			return ""
		}
		return b2s(input[:pos1])
	}
	var i uint8
	for i = 0; i < segment-1; i++ {
		input = input[pos1:]
		pos1 = bytes.IndexByte(input, '/')
		if pos1 < 0 {
			return ""
		}
		pos1++
	}
	input = input[pos1:]
	pos2 := bytes.IndexByte(input, '/')
	if pos2 < 0 {
		return b2s(input)
	}
	return b2s(input[:pos2])
}

/*
func (node *Node) GetParam(pattern, name string, w io.Writer) {
	curr := node
	for curr != nil {
		if curr.code == ':' && curr.param == name {
			matchParam(curr, pattern, name, w)
			return
		}
		curr = curr.parent
	}
}

func matchParam(node *Node, pattern, name string, w io.Writer) {
	input := unsafe.Slice(unsafe.StringData(pattern), len(pattern))

	chain := nodeSlicePool.Get().(*NodeSlice)
	defer nodeSlicePool.Put(chain)
	chain.items = chain.items[:0]

	curr := node
	for curr != nil {
		chain.items = append(chain.items, curr)
		curr = curr.parent
	}
	slices.Reverse(chain.items)
	offset := 0
	skip := false
	catch := false

	for i, c1 := range input {

		if skip {
			if c1 != '/' {
				offset++
				continue
			} else {
				skip = false
			}
		}

		if catch {
			if c1 == '/' {
				break
			}
			w.Write([]byte{c1})
			continue
		}

		c2 := chain.items[i-offset]

		if c1 == c2.code {
			continue
		}

		if c2.code == ':' {
			if c2.param != name {
				skip = true
				continue
			} else {
				w.Write([]byte{c1})
				catch = true
			}
		}
	}
}
*/
