package httpx

import (
	"net/http"

	"github.com/google/uuid"
)

const TraceIDHeader = "X-Trace-Id"

func TraceID(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if w.Header().Get(TraceIDHeader) == "" {
			w.Header().Set(TraceIDHeader, GetOrCreateTraceID(r))
		}
		return next(w, r)
	}
}

func GetOrCreateTraceID(r *http.Request) string {
	traceID := GetTraceID(r)
	if traceID == "" {
		traceID = uuid.New().String()
		r.Header.Set(TraceIDHeader, traceID)
	}
	return traceID
}

func GetTraceID(r *http.Request) string {
	return r.Header.Get(TraceIDHeader)
}
