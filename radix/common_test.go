package radix

import (
	"bytes"
	"testing"
)

func TestExtractParam(t *testing.T) {
	testCases := [][3][]byte{
		{[]byte(":param1"), []byte("param1"), nil},
		{[]byte("param1"), []byte("param1"), nil},
		{[]byte(":param1/other"), []byte("param1"), []byte("/other")},
		{[]byte("param1/other"), []byte("param1"), []byte("/other")},
	}
	for _, tc := range testCases {
		param, tail := extractParam([]byte(tc[0]))
		if !bytes.Equal(param, tc[1]) {
			t.Fatalf("expect param '%s', actual '%s'", string(tc[1]), string(param))
		}
		if !bytes.Equal(tail, tc[2]) {
			t.Fatalf("expect tail '%s', actual '%s'", string(tc[2]), string(tail))
		}
	}
}

func TestGetSegment(t *testing.T) {
	input := []byte("foo/bar/baz")
	seg := getSegment(input, 1)
	if seg != "foo" {
		t.Fatalf("wrong segment, expect 'foo', actual '%s'", seg)
	}
	seg = getSegment(input, 2)
	if seg != "bar" {
		t.Fatalf("wrong segment, expect 'bar', actual '%s'", seg)
	}
	seg = getSegment(input, 3)
	if seg != "baz" {
		t.Fatalf("wrong segment, expect 'baz', actual '%s'", seg)
	}

}
