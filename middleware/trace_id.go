package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nskforward/httpx"
)

func TraceID(req *http.Request, resp *httpx.Response) error {
	traceID := req.Header.Get("X-Trace-Id")
	if traceID == "" {
		traceID = uuid.New().String()
		req.Header.Set("X-Trace-Id", traceID)
	}
	resp.LoggingWith("trace-id", traceID)
	return resp.Next(req)
}
