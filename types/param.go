package types

import (
	"context"
	"net/http"
)

type ContextParam string

func SetParam(r *http.Request, param ContextParam, value any) *http.Request {
	ctx := context.WithValue(r.Context(), param, value)
	return r.WithContext(ctx)
}

func GetParam(r *http.Request, param ContextParam) any {
	return r.Context().Value(param)
}
