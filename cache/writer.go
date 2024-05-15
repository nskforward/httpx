package cache

import (
	"bytes"
	"net/http"
	"time"
)

type Writer struct {
	code    int
	size    int
	maxSize int
	age     time.Duration
	raw     http.ResponseWriter
	ignore  bool
	buf     bytes.Buffer
	reason  string
}

func (w *Writer) Reset(raw http.ResponseWriter, maxSize int) {
	w.code = 0
	w.size = 0
	w.maxSize = maxSize
	w.age = 0
	w.raw = raw
	w.ignore = false
	w.buf.Reset()
	w.reason = ""
}

func (w *Writer) StatusCode() int {
	return w.code
}

func (w *Writer) CacheAge() time.Duration {
	return w.age
}

func (w *Writer) Buffer() *bytes.Buffer {
	return &w.buf
}

func (w *Writer) Header() http.Header {
	return w.raw.Header()
}

func (w *Writer) CanCache() bool {
	return !w.ignore
}

func (w *Writer) Reason() string {
	return w.reason
}

func (w *Writer) checkHeaders(status int) {
	w.code = status

	if status != 200 {
		// don't cache any responses except 200 OK
		w.ignore = true
		w.reason = "bad status"
		return
	}

	cacheControl := w.Header().Get("Cache-Control")
	if cacheControl == "" {
		w.ignore = true
		w.reason = "no control"
		return
	}

	cc := NewControl(cacheControl)
	if !cc.IsPublic {
		w.ignore = true
		w.reason = "not public"
		return
	}

	if cc.MaxAge > 0 {
		w.age = cc.MaxAge
	}
}

func (w *Writer) WriteHeader(status int) {
	w.checkHeaders(status)
	w.raw.WriteHeader(w.code)
}

func (w *Writer) Write(b []byte) (int, error) {
	if w.code == 0 {
		w.WriteHeader(200)
	}

	n, err := w.raw.Write(b)

	if w.ignore {
		return n, err
	}

	// try to save cache
	if w.buf.Len()+len(b) > w.maxSize {
		w.ignore = true
		w.reason = "too large"
		return n, err
	}

	return w.buf.Write(b)
}
