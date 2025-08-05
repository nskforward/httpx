package httpx

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
	status              int
	size                int64
	onBeforeHeadersSent []func(resp *Response)
}

func newResponse(w http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter:      w,
		status:              0,
		onBeforeHeadersSent: make([]func(resp *Response), 8),
	}
}

func (ww *Response) OnBeforeHeadersSent(subscruber func(*Response)) {
	ww.onBeforeHeadersSent = append(ww.onBeforeHeadersSent, subscruber)
}

func (ww *Response) Size() int64 {
	return ww.size
}

func (ww *Response) Status() int {
	return ww.status
}

func (ww *Response) HeadersSent() bool {
	return ww.status > 0
}

func (ww *Response) WriteHeader(statusCode int) {
	ww.status = statusCode
	for _, subscriber := range ww.onBeforeHeadersSent {
		if subscriber != nil {
			subscriber(ww)
		}
	}
	ww.ResponseWriter.WriteHeader(statusCode)
}

func (ww *Response) Write(p []byte) (int, error) {
	if ww.status == 0 {
		ww.WriteHeader(200)
	}
	written, err := ww.ResponseWriter.Write(p)
	ww.size += int64(written)
	return written, err
}

func (ww *Response) Flusher() http.Flusher {
	flusher, ok := ww.ResponseWriter.(http.Flusher)
	if !ok {
		panic("w http.ResponseWriter does not implement http.Flusher")
	}
	return flusher
}

func (ww *Response) SendText(code int, msg string) error {
	http.Error(ww.ResponseWriter, msg, code)
	return nil
}

func (ww *Response) SendJSON(code int, object any) error {
	ww.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(ww).Encode(object)
}
