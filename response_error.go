package httpx

import (
	"fmt"
	"net/http"
	"strings"
)

type ResponseError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	ContentType string `json:"content-type,omitempty"`
	Source      string `json:"source,omitempty"`
}

func NewError(code int, err error) ResponseError {
	return ResponseError{
		Code:    code,
		Message: err.Error(),
	}
}

func (e ResponseError) Error() string {
	return e.Message
}

func (e ResponseError) WithJSON() ResponseError {
	e.ContentType = ContentTypeJSON
	return e
}

func (e ResponseError) WithSource(source string) ResponseError {
	e.Source = source
	return e
}

func (e ResponseError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e.Code == 0 {
		e.Code = http.StatusBadRequest
	}

	contentType := e.ContentType
	if contentType == "" {
		contentType = detectContentType(r, ContentTypeTEXT)
	}
	w.Header().Set("Content-Type", contentType)
	if strings.Contains(contentType, "application/json") {
		JSON(w, e.Code, e)
		return
	}
	Text(w, e.Code, e.Message)
}

func BadGateway(err error) ResponseError {
	return NewError(http.StatusBadGateway, err)
}

func BadRequest(err error) ResponseError {
	return NewError(http.StatusBadRequest, err)
}

func InternalServerError(err error) ResponseError {
	return NewError(http.StatusInternalServerError, err)
}

func Unauthorized(err error) ResponseError {
	return NewError(http.StatusUnauthorized, err)
}

func Forbidden(err error) ResponseError {
	return NewError(http.StatusForbidden, err)
}

func NotFound() ResponseError {
	return NewError(http.StatusNotFound, fmt.Errorf(http.StatusText(http.StatusNotFound)))
}

func MethodNotAllowed() ResponseError {
	return NewError(http.StatusMethodNotAllowed, fmt.Errorf(http.StatusText(http.StatusMethodNotAllowed)))
}

func detectContentType(r *http.Request, defaultContentType string) string {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return defaultContentType
	}
	if strings.Contains(contentType, "application/json") {
		return ContentTypeJSON
	}
	return defaultContentType
}
