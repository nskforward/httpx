package httpx

import (
	"io"
	"net/http"
)

type responseWrapper struct {
	http.ResponseWriter
	status int
	size   int64
	body   io.Writer
}

func newResponseWrapper(w http.ResponseWriter) *responseWrapper {
	return &responseWrapper{ResponseWriter: w, status: 0, body: w}
}

func (ww *responseWrapper) Size() int64 {
	return ww.size
}

func (ww *responseWrapper) Status() int {
	return ww.status
}

func (ww *responseWrapper) Header() http.Header {
	return ww.ResponseWriter.Header()
}

func (ww *responseWrapper) WriteHeader(statusCode int) {
	ww.status = statusCode
	ww.ResponseWriter.WriteHeader(statusCode)
}

func (ww *responseWrapper) Write(p []byte) (int, error) {
	if ww.status == 0 {
		ww.WriteHeader(200)
	}
	written, err := ww.body.Write(p)
	ww.size += int64(written)
	return written, err
}

func (ww *responseWrapper) Flusher() http.Flusher {
	flusher, ok := ww.ResponseWriter.(http.Flusher)
	if !ok {
		panic("w http.ResponseWriter does not implement http.Flusher")
	}
	return flusher
}
