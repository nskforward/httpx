package httpx

import (
	"net/http"

	"github.com/google/uuid"
)

const TraceIDHeader = "X-Trace-Id"

func NewTraceID(r *http.Request) string {
	traceID := r.Header.Get(TraceIDHeader)
	if traceID == "" {
		traceID = uuid.New().String()
	}
	return traceID
}
