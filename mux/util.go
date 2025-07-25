package mux

import (
	"bytes"
	"fmt"
	"iter"
	"strings"
)

func SplitPattern(pattern string) (methodNum Method, url string, err error) {
	method := "ANY"

	if len(pattern) == 0 {
		err = fmt.Errorf("pattern cannot be empty")
		return
	}

	if pattern[0] != '/' {
		found := false
		method, url, found = strings.Cut(pattern, " ")
		if !found {
			err = fmt.Errorf("bad pattern format")
			return
		}
	} else {
		url = pattern
	}

	if url[0] != '/' {
		err = fmt.Errorf("url must start with '/'")
	}

	methodNum = MethodToUInt8(method)

	return
}

func Segments(url string) iter.Seq2[string, bool] {
	var buf bytes.Buffer
	size := 0
	return func(yield func(string, bool) bool) {
		for _, char := range url {
			size++
			if char == '/' {
				if buf.Len() > 0 {
					if !yield(buf.String(), size == len(url)) {
						return
					}
					buf.Reset()
				}
				if !yield("/", size == len(url)) {
					return
				}
				continue
			}
			buf.WriteRune(char)
		}
		if buf.Len() > 0 {
			if !yield(buf.String(), size == len(url)) {
				return
			}
			buf.Reset()
		}
	}
}
