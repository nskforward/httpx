package httpx

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type Response struct {
	logger   *slog.Logger
	w        *ResponseWrapper
	handlers []Handler
	index    int
}

func NewResponse(logger *slog.Logger, w http.ResponseWriter, handler []Handler) *Response {
	if logger == nil {
		panic("httpx.NewResponse requres not nil logger")
	}
	return &Response{
		w:        NewResponseWrapper(w),
		handlers: handler,
		index:    0,
		logger:   logger,
	}
}

func (resp *Response) ResponseWriter() http.ResponseWriter {
	return resp.w
}

func (resp *Response) StatusCode() int {
	return resp.w.status
}

func (resp *Response) BodySize() int64 {
	return resp.w.size
}

func (resp *Response) Logger() *slog.Logger {
	return resp.logger
}

func (resp *Response) SetHeader(name, value string) {
	resp.w.Header().Set(name, value)
}

func (resp *Response) LoggingWith(args ...any) {
	resp.logger = resp.logger.With(args...)
}

func (resp *Response) InternalServerError() error {
	return resp.Text(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (resp *Response) Unauthorized() error {
	return resp.Text(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}

func (resp *Response) Next(req *http.Request) error {
	if resp.index < len(resp.handlers) {
		next := resp.handlers[resp.index]
		resp.index++
		return next(req, resp)
	}
	return nil
}

func (resp *Response) Text(code int, msg string) error {
	resp.w.WriteHeader(code)
	io.WriteString(resp.w, msg)
	return nil
}

func (resp *Response) Write(code int, src []byte) error {
	resp.w.WriteHeader(code)
	resp.w.Write(src)
	return nil
}

func (resp *Response) Copy(code int, src io.Reader) error {
	resp.w.WriteHeader(code)
	io.Copy(resp.w, src)
	return nil
}

func (resp *Response) JSON(code int, obj any) error {
	resp.SetHeader("Content-Type", "application/json; charset=utf-8")
	resp.w.WriteHeader(code)
	return json.NewEncoder(resp.w).Encode(obj)
}

func (resp *Response) NoContent() error {
	resp.w.WriteHeader(http.StatusNoContent)
	return nil
}
