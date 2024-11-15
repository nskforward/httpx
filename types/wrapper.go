package types

import (
	"io"
	"net/http"
	"time"
)

type ResponseWrapper struct {
	http.ResponseWriter
	http.Flusher
	status      int
	size        int64
	BeforeBody  func()
	AllowHeader func(name string, values []string) bool
	body        io.Writer
	wroteHeader bool
	started     time.Time
	skipBody    bool
	writing     bool
}

func NewResponseWrapper(w http.ResponseWriter) *ResponseWrapper {
	return &ResponseWrapper{ResponseWriter: w, status: 404, body: w, started: time.Now()}
}

func (ww *ResponseWrapper) Size() int64 {
	return ww.size
}

func (ww *ResponseWrapper) StartTime() time.Time {
	return ww.started
}

func (ww *ResponseWrapper) Status() int {
	return ww.status
}

func (ww *ResponseWrapper) SkipBody() {
	ww.skipBody = true
}

func (ww *ResponseWrapper) SetWriter(w io.Writer) {
	ww.body = w
}

func (ww *ResponseWrapper) WriteHeader(statusCode int) {
	if ww.wroteHeader {
		return
	}

	ww.status = statusCode

	if ww.AllowHeader != nil {
		for name, values := range ww.Header() {
			if !ww.AllowHeader(name, values) {
				ww.Header().Del(name)
			}
		}
	}

	if ww.BeforeBody != nil && !ww.writing {
		ww.writing = true
		ww.BeforeBody()
		ww.writing = false
	}
	if ww.wroteHeader {
		return
	}
	if ww.body == nil {
		panic("response.Writer body is nil")
	}
	ww.ResponseWriter.WriteHeader(statusCode)
	ww.wroteHeader = true
}

func (ww *ResponseWrapper) Write(p []byte) (written int, err error) {
	if !ww.wroteHeader {
		ww.WriteHeader(200)
	}
	if ww.skipBody {
		return len(p), nil
	}
	written, err = ww.body.Write(p)
	ww.size += int64(written)
	return
}

func (ww *ResponseWrapper) Flush() {
	flusher, ok := ww.ResponseWriter.(http.Flusher)
	if !ok {
		panic("w http.ResponseWriter does not implement http.Flusher")
	}
	flusher.Flush()
}
