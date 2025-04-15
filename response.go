package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type Response struct {
	logger *slog.Logger
	w      *ResponseWrapper
	index  int
	h      []Handler
}

func NewResponse(logger *slog.Logger, w http.ResponseWriter, h []Handler) *Response {
	if logger == nil {
		panic("httpx.NewResponse requres not nil logger")
	}
	return &Response{
		w:      NewResponseWrapper(w),
		logger: logger,
		h:      h,
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

func (resp *Response) SetContentType(contentType string) {
	resp.SetHeader("Content-Type", contentType)
}

func (resp *Response) LoggingWith(args ...any) {
	resp.logger = resp.logger.With(args...)
}

func (resp *Response) ServerError(err error) error {
	resp.logger.Error("internal server error", "error", err.Error())
	return resp.Text(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (resp *Response) ClientError(err error) error {
	apiError, ok := err.(*APIError)
	if ok {
		return resp.Text(apiError.Code, apiError.Mesage)
	}
	return resp.Text(http.StatusBadRequest, err.Error())
}

func (resp *Response) Unauthorized(msg string) error {
	return resp.Text(http.StatusUnauthorized, msg)
}

func (resp *Response) Forbidden(msg string) error {
	return resp.Text(http.StatusForbidden, msg)
}

func (resp *Response) Next(req *http.Request) error {
	if resp.index < len(resp.h) {
		resp.index++
		return resp.h[resp.index-1](req, resp)
	}
	return errors.New("no next handler found")
}

func (resp *Response) Text(code int, msg string) error {
	resp.SetContentType("text/plain; charset=UTF-8")
	resp.w.WriteHeader(code)
	io.WriteString(resp.w, msg)
	return nil
}

func (resp *Response) JSON(code int, obj any) error {
	resp.SetContentType("application/json; charset=utf-8")
	resp.w.WriteHeader(code)
	return json.NewEncoder(resp.w).Encode(obj)
}

func (resp *Response) WriteBytes(code int, contentType string, src []byte) error {
	resp.SetHeader("Content-Type", contentType)
	resp.w.WriteHeader(code)
	resp.w.Write(src)
	return nil
}

func (resp *Response) WriteReader(code int, contentType string, src io.Reader) error {
	resp.SetHeader("Content-Type", contentType)
	resp.w.WriteHeader(code)
	io.Copy(resp.w, src)
	return nil
}

func (resp *Response) NoContent() error {
	resp.w.WriteHeader(http.StatusNoContent)
	return nil
}
