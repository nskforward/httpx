package radix

import (
	"fmt"
	"math/rand"
	"testing"
)

type testCase struct {
	pattern string
	path    string
	params  map[string]string
	value   int
}

var testCases = []testCase{
	{
		pattern: "/foo",
		path:    "/foo",
		value:   1,
	},
	{
		pattern: "/",
		path:    "/",
		value:   2,
	},
	{
		pattern: "/foo/1",
		path:    "/foo/1",
		value:   3,
	},
	{
		pattern: "/foo/1234",
		path:    "/foo/1234",
		value:   4,
	},
	{
		pattern: "/bar/1",
		path:    "/bar/1",
		value:   5,
	},
	{
		pattern: "/foo/bar/1",
		path:    "/foo/bar/1",
		value:   6,
	},
	{
		pattern: "/*",
		path:    "/foo/bar/1/baz",
		value:   7,
	},
	{
		pattern: "/foo/:foo",
		path:    "/foo/2",
		params:  map[string]string{"foo": "2"},
		value:   8,
	},
	{
		pattern: "/foo/bar/:bar",
		path:    "/foo/bar/2",
		params:  map[string]string{"bar": "2"},
		value:   9,
	},
	{
		pattern: "/foo/:foo/bar/:bar/:baz",
		path:    "/foo/1/bar/2/3",
		params:  map[string]string{"foo": "1", "bar": "2", "baz": "3"},
		value:   10,
	},
}

func createMockMux() *Node {
	var mux Node
	for _, tc := range testCases {
		err := mux.Insert(tc.pattern, tc.value)
		if err != nil {
			panic(fmt.Errorf("insert error: '%s': %w", tc.pattern, err))
		}
	}
	return &mux
}

func TestNodeSearch(t *testing.T) {
	mux := createMockMux()
	DumpTree(mux, 0)
	for _, tc := range testCases {
		node := mux.Search(tc.path)
		if node == nil {
			t.Fatalf("cannot find node on path: %s", tc.path)
		}
		if node.value != tc.value {
			t.Fatalf("wrong value, path=%s, wants=%d, actual=%d", tc.path, tc.value, node.value)
		}
		for k, v1 := range tc.params {
			v2 := node.GetParam(tc.path, k)
			if v1 != v2 {
				t.Fatalf("wrong param, path=%s, name=%s wants=%s, actual=%s", tc.path, k, v1, v2)
			}
		}
	}
}

func BenchmarkNodeSearch(b *testing.B) {
	mux := createMockMux()
	for i := 0; i < b.N; i++ {
		index := rand.Intn(len(testCases))
		tc := testCases[index]
		node := mux.Search(tc.path)
		if node == nil {
			b.Fatalf("cannot find node on path: %s", tc.path)
		}
		if node.value != tc.value {
			b.Fatalf("wrong value, path=%s, wants=%d, actual=%d", tc.path, index, node.value)
		}
		for k, v1 := range tc.params {
			v2 := node.GetParam(tc.path, k)
			if v1 != v2 {
				b.Fatalf("wrong param, path=%s, name=%s wants=%s, actual=%s", tc.path, k, v1, v2)
			}
		}
	}
}
