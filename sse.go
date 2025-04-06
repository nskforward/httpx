package httpx

import (
	"io"
	"net/http"
)

type SSE struct {
	w          http.ResponseWriter
	controller *http.ResponseController
}

func (resp *Response) Stream() *SSE {
	resp.SetHeader("Content-Type", "text/event-stream")
	resp.SetHeader("Cache-Control", "no-store")

	return &SSE{
		w:          resp.w,
		controller: http.NewResponseController(resp.w),
	}
}

func (sse *SSE) Flush() error {
	_, err := sse.w.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return sse.controller.Flush()
}

func (sse *SSE) WriteField(name, value string) error {
	_, err := io.WriteString(sse.w, name)
	if err != nil {
		return err
	}
	_, err = sse.w.Write([]byte{':', ' '})
	if err != nil {
		return err
	}
	_, err = io.WriteString(sse.w, value)
	if err != nil {
		return err
	}
	_, err = sse.w.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return nil
}
