package types

import (
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) error
type Middleware func(next Handler) Handler
type LoggerFunc func(w *ResponseWrapper, r *http.Request, err error)
