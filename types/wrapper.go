package types

import (
	"io"
	"net/http"
	"time"
)

type ResponseWrapper struct {
	http.ResponseWriter
	status      int
	size        int64
	BeforeBody  func()
	body        io.Writer
	wroteHeader bool
	started     time.Time
}

func NewResponseWrapper(w http.ResponseWriter) *ResponseWrapper {
	return &ResponseWrapper{ResponseWriter: w, status: 200, body: w, started: time.Now()}
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

func (ww *ResponseWrapper) SetWriter(w io.Writer) {
	ww.body = w
}

func (ww *ResponseWrapper) WriteHeader(statusCode int) {
	if ww.wroteHeader {
		return
	}
	ww.status = statusCode
	if ww.BeforeBody != nil {
		ww.BeforeBody()
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
	written, err = ww.body.Write(p)
	ww.size += int64(written)
	return
}
