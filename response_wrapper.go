package httpx

import (
	"io"
	"net/http"
)

type ResponseWrapper struct {
	http.ResponseWriter
	status int
	size   int64
	body   io.Writer
}

func NewResponseWrapper(w http.ResponseWriter) *ResponseWrapper {
	return &ResponseWrapper{ResponseWriter: w, status: 200, body: w}
}

func (ww *ResponseWrapper) Size() int64 {
	return ww.size
}

func (ww *ResponseWrapper) Status() int {
	return ww.status
}

func (ww *ResponseWrapper) Header() http.Header {
	return ww.ResponseWriter.Header()
}

func (ww *ResponseWrapper) WriteHeader(statusCode int) {
	ww.status = statusCode
	ww.ResponseWriter.WriteHeader(statusCode)
}

func (ww *ResponseWrapper) Write(p []byte) (int, error) {
	written, err := ww.body.Write(p)
	ww.size += int64(written)
	return written, err
}

func (ww *ResponseWrapper) Flusher() http.Flusher {
	flusher, ok := ww.ResponseWriter.(http.Flusher)
	if !ok {
		panic("w http.ResponseWriter does not implement http.Flusher")
	}
	return flusher
}
