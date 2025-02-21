package httpx

import (
	"fmt"
	"net/http"
)

type Error struct {
	code int
	text string
}

var (
	ErrNotFound        = Error{code: http.StatusNotFound, text: http.StatusText(http.StatusNotFound)}
	ErrInternalServer  = Error{code: http.StatusInternalServerError, text: http.StatusText(http.StatusInternalServerError)}
	ErrUnauthorized    = Error{code: http.StatusUnauthorized, text: http.StatusText(http.StatusUnauthorized)}
	ErrForbidden       = Error{code: http.StatusForbidden, text: http.StatusText(http.StatusForbidden)}
	ErrTooManyRequests = Error{code: http.StatusTooManyRequests, text: http.StatusText(http.StatusTooManyRequests)}
)

func NewError(code int, text string) Error {
	return Error{
		code: code,
		text: text,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.text)
}

func (e Error) Write(ctx *Ctx) {
	ctx.Text(e.code, e.text)
}
