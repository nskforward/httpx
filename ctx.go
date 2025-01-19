package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Context struct {
	logger      *slog.Logger
	w           *ResponseWrapper
	req         *http.Request
	pattern     string
	traceID     string
	headersSent bool
	startTime   time.Time
}

func NewContext(parent *slog.Logger, pattern string, w http.ResponseWriter, req *http.Request) *Context {
	traceID := NewTraceID(req)

	return &Context{
		logger: parent.With(
			slog.String("trace_id", traceID),
			slog.String("pattern", pattern),
		),
		w:         NewResponseWrapper(w),
		req:       req,
		pattern:   pattern,
		traceID:   traceID,
		startTime: time.Now(),
	}
}

func (ctx *Context) ParseRequestJSON(dest any) error {
	return json.NewDecoder(ctx.req.Body).Decode(dest)
}

func (ctx *Context) Request() *http.Request {
	return ctx.req
}

func (ctx *Context) PathParam(name string) string {
	return ctx.req.PathValue(name)
}

func (ctx *Context) FormParam(name string) string {
	return ctx.req.FormValue(name)
}

func (ctx *Context) TraceID() string {
	return ctx.traceID
}

func (ctx *Context) Method() string {
	return ctx.req.Method
}

func (ctx *Context) UserAgent() string {
	return ctx.req.UserAgent()
}

func (ctx *Context) UserIP() string {
	ip, _, _ := net.SplitHostPort(ctx.req.RemoteAddr)
	return ip
}

func (ctx *Context) Path() string {
	return ctx.req.URL.Path
}

func (ctx *Context) StatusCode() int {
	return ctx.w.status
}

func (ctx *Context) ClientContext() context.Context {
	return ctx.req.Context()
}

func (ctx *Context) StartTime() time.Time {
	return ctx.startTime
}

func (ctx *Context) Logger() *slog.Logger {
	return ctx.logger
}

func (ctx *Context) ResponseSize() int64 {
	return ctx.w.size
}

func (ctx *Context) HeadersSent() bool {
	return ctx.headersSent
}

func (ctx *Context) sendStatusCode(statusCode int) {
	ctx.headersSent = true
	ctx.w.WriteHeader(statusCode)
}

func (ctx *Context) GetRequestHeader(name string) string {
	return ctx.req.Header.Get(name)
}

func (ctx *Context) SetResponseHeader(name, value string) {
	if ctx.headersSent {
		ctx.Logger().Warn("try to set an http response header when status code already sent to client", "name", name, "value", value)
		return
	}
	ctx.w.Header().Set(name, value)
}

func (ctx *Context) RespondNoContent() error {
	ctx.sendStatusCode(http.StatusNoContent)
	return nil
}

func (ctx *Context) RespondJSON(statusCode int, obj any) error {
	ctx.SetResponseHeader("Content-Type", "application/json; charset=utf-8")
	ctx.sendStatusCode(statusCode)
	return json.NewEncoder(ctx.w).Encode(obj)
}

func (ctx *Context) RespondText(statusCode int, msg string) error {
	ctx.sendStatusCode(statusCode)
	_, err := io.WriteString(ctx.w, msg)
	return err
}

func (ctx *Context) Redirect(statusCode int, url string) error {
	http.Redirect(ctx.w, ctx.req, url, statusCode)
	return nil
}

// Stream returns TRUE if client gone, FALSE if server breaks stream.
func (ctx *Context) Stream(step func(send func(name, value string), flush func()) bool) bool {
	ctx.SetResponseHeader("Content-Type", "text/event-stream")
	ctx.SetResponseHeader("Cache-Control", "no-store")

	gone := ctx.req.Context().Done()
	writer := ctx.w
	flusher := ctx.w.Flusher()

	send := func(name, value string) {
		writer.Write([]byte(name))
		writer.Write([]byte(": "))
		writer.Write([]byte(value))
		writer.Write([]byte{'\n'})
	}

	flush := func() {
		writer.Write([]byte{'\n'})
		flusher.Flush()
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-gone:
			return true

		case <-ticker.C:
			writer.Write([]byte(":ping\n"))
			flusher.Flush()

		default:
			keepOpen := step(send, flush)
			if !keepOpen {
				return false
			}
		}
	}
}

func (ctx *Context) AccessDenied() error {
	return ctx.RespondText(http.StatusForbidden, "access denied")
}

func (ctx *Context) BadRequest(msg string) error {
	return ctx.RespondText(http.StatusBadRequest, msg)
}

func (ctx *Context) CacheDisable() {
	ctx.SetResponseHeader("Cache-Control", "no-store")
}

func (ctx *Context) CacheEnable(public bool, maxAge time.Duration) {
	ctx.SetResponseHeader("Cache-Control", fmt.Sprintf("%s, max-age=%.0f", isPublic(public), maxAge.Seconds()))
}

func isPublic(yes bool) string {
	if yes {
		return "public"
	}
	return "private"
}
