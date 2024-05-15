package logging

import (
	"net/http"
	"time"
)

type Writer struct {
	code  int
	size  int
	raw   http.ResponseWriter
	start time.Time
}

func (w *Writer) Reset(raw http.ResponseWriter) {
	w.code = 0
	w.size = 0
	w.raw = raw
	w.start = time.Now()
}

func (w *Writer) StatusCode() int {
	return w.code
}

func (w *Writer) Duration() time.Duration {
	return time.Since(w.start)
}

func (w *Writer) Size() int {
	return w.size
}

func (w *Writer) Header() http.Header {
	return w.raw.Header()
}

func (w *Writer) WriteHeader(status int) {
	w.code = status
	w.raw.WriteHeader(status)
}

func (w *Writer) Write(data []byte) (int, error) {
	if w.code == 0 {
		w.WriteHeader(200)
	}
	w.size += len(data)
	return w.raw.Write(data)
}

func (w *Writer) Close() {
	if f, ok := w.raw.(http.Flusher); ok {
		f.Flush()
	}
}
