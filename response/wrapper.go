package response

import (
	"io"
	"net/http"
)

type Wrapper struct {
	http.ResponseWriter
	status      int
	size        int64
	BeforeBody  func()
	body        io.Writer
	headersSent bool
}

func NewWrapper(w http.ResponseWriter) *Wrapper {
	return &Wrapper{ResponseWriter: w, status: 200, size: 0, body: w, headersSent: false}
}

func (ww *Wrapper) Size() int64 {
	return ww.size
}

func (ww *Wrapper) Status() int {
	return ww.status
}

func (ww *Wrapper) BodyWriter(w io.Writer) {
	ww.body = w
}

func (ww *Wrapper) WriteHeader(statusCode int) {
	ww.status = statusCode
	if ww.BeforeBody != nil {
		ww.BeforeBody()
	}
	ww.ResponseWriter.WriteHeader(statusCode)
	ww.headersSent = true
}

func (ww *Wrapper) Write(p []byte) (written int, err error) {
	if !ww.headersSent {
		ww.WriteHeader(200)
	}
	written, err = ww.body.Write(p)
	ww.size += int64(written)
	return
}
