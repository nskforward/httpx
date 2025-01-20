package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

type Context struct {
	logger      *slog.Logger
	w           *ResponseWrapper
	req         *http.Request
	traceID     string
	headersSent bool
	startTime   time.Time
	realIP      string
}

func NewContext(parent *slog.Logger, w http.ResponseWriter, req *http.Request) *Context {
	traceID := NewTraceID(req)

	return &Context{
		logger: parent.With(
			slog.String("trace_id", traceID),
		),
		w:         NewResponseWrapper(w),
		req:       req,
		traceID:   traceID,
		startTime: time.Now(),
		realIP:    detectRealIP(req),
	}
}

func (ctx *Context) ParseRequestJSON(dest any) error {
	return json.NewDecoder(ctx.req.Body).Decode(dest)
}

func (ctx *Context) Request() *http.Request {
	return ctx.req
}

func (ctx *Context) ResponseWriter() http.ResponseWriter {
	return ctx.w
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
	return ctx.realIP
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

func (ctx *Context) RespondRestResponse(success bool, description string, payload any) error {
	return ctx.RespondJSON(200, RestResponse{Success: success, Description: description, Payload: payload})
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

func (ctx *Context) Unauthorized() error {
	return ctx.RespondText(http.StatusUnauthorized, "unauthorized")
}

func (ctx *Context) BadRequest(msg string, args ...any) error {
	args = append([]any{msg}, args...)
	args = append(args, fmt.Sprintf("(#%s)", ctx.TraceID()))
	return ctx.RespondText(http.StatusBadRequest, fmt.Sprint(args...))
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

func detectRealIP(r *http.Request) string {
	addr := r.Header.Get("X-Forwarded-For")
	if addr != "" {
		if strings.Contains(addr, ",") {
			addr = strings.TrimSpace(strings.Split(addr, ",")[0])
		}
		return addr
	}
	addr = r.Header.Get("X-Real-IP")
	if addr != "" {
		return addr
	}
	addr, _, _ = net.SplitHostPort(r.RemoteAddr)
	return addr
}
