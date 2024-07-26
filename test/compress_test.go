package test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/middleware"
	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func TestCompress(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println("--- main handler")
		return response.Text(w, 200, strings.Repeat("0123456789", 201))
	}

	var r httpx.Router
	r.Use(middleware.Compress)
	r.Route("/api/v1/", h)

	s := httptest.NewServer(&r)
	defer s.Close()

	DoRequest(s, "GET", "/api/v1/user/123", "", http.Header{types.AcceptEncoding: []string{"gzip"}}, true, false)
}

var testData = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur molestie enim a urna mattis, eu congue ipsum fermentum. Aliquam ullamcorper luctus viverra. Nullam id purus magna. Duis sed cursus metus. Sed vitae risus laoreet, volutpat eros ut, posuere dui. Vestibulum tellus mi, vestibulum eget mauris in, volutpat auctor dolor. Maecenas in auctor libero. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur molestie enim a urna mattis, eu congue ipsum fermentum. Aliquam ullamcorper luctus viverra. Nullam id purus magna. Duis sed cursus metus. Sed vitae risus laoreet, volutpat eros ut, posuere dui. Vestibulum tellus mi, vestibulum eget mauris in, volutpat auctor dolor. Maecenas in auctor libero.")

func TestGZip(t *testing.T) {
	var buf bytes.Buffer
	err := compressTestData(&buf, testData, getGzWriterFresh, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("reduction:", 100-buf.Len()*100/len(testData), "%")
}

/*
BenchmarkGZipPoolYes-4(buf pool no)		1917	    542099 ns/op	  207807 B/op	       8 allocs/op
BenchmarkGZipPoolYes-4(buf pool yes)	2000	    539329 ns/op	  206944 B/op	       5 allocs/op
BenchmarkGZipPoolNo-4	1206	    834749 ns/op	  813880 B/op	      17 allocs/op
*/
func BenchmarkGZipPoolYes(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := getBuffPool()
			err := compressTestData(buf, testData, getGzWriterPool, putGzWriterPool)
			putBuffPool(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGZipPoolNo(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := getBuffFresh()
			err := compressTestData(buf, testData, getGzWriterFresh, nil)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

var gzPool = sync.Pool{
	New: func() any {
		return gzip.NewWriter(io.Discard)
	},
}

var bytePool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func compressTestData(w io.Writer, testData []byte, getWriter func(io.Writer) *gzip.Writer, close func(*gzip.Writer)) error {
	gz := getWriter(w)
	_, err := gz.Write(testData)
	gz.Close()
	if close != nil {
		close(gz)
	}
	return err
}

func getGzWriterFresh(w io.Writer) *gzip.Writer {
	return gzip.NewWriter(w)
}

func getGzWriterPool(w io.Writer) *gzip.Writer {
	gz := gzPool.Get().(*gzip.Writer)
	gz.Reset(w)
	return gz
}

func putGzWriterPool(w *gzip.Writer) {
	gzPool.Put(w)
}

func getBuffFresh() *bytes.Buffer {
	return new(bytes.Buffer)
}

func getBuffPool() *bytes.Buffer {
	buf := bytePool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func putBuffPool(buf *bytes.Buffer) {
	bytePool.Put(buf)
}
