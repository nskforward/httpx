package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Ctx struct {
	r            *http.Request
	w            *responseWrapper
	route        *Route
	indexHandler int
	clientAddr   string
	traceID      string
	logger       *slog.Logger
	formParsed   bool
}

func newCtx(route *Route, w http.ResponseWriter, r *http.Request) *Ctx {
	traceID := NewTraceID(r)
	return &Ctx{
		r:            r,
		w:            newResponseWrapper(w),
		route:        route,
		indexHandler: -1,
		clientAddr:   r.RemoteAddr,
		traceID:      traceID,
		logger:       route.router.app.logger.With(slog.String("trace_id", traceID)),
	}
}

func (ctx *Ctx) TraceID() string {
	return ctx.traceID
}

func (ctx *Ctx) Context() context.Context {
	return ctx.Request().Context()
}

func (ctx *Ctx) Path() string {
	return ctx.r.URL.Path
}

func (ctx *Ctx) Origin() string {
	return ctx.r.Header.Get("Origin")
}

func (ctx *Ctx) Request() *http.Request {
	return ctx.r
}

func (ctx *Ctx) Logger() *slog.Logger {
	return ctx.logger
}

func (ctx *Ctx) Next() error {
	nextHandler := ctx.nextHandler()
	if nextHandler != nil {
		return nextHandler(ctx)
	}
	return ErrNotFound
}

func (ctx *Ctx) Sent() bool {
	return ctx.w.Status() > 0
}

func (ctx *Ctx) WriteError(err error) {
	if ctx.Sent() {
		ctx.Logger().Warn("response headers already sent", "error", err)
		return
	}

	fmt.Println("catch error")

	if servErr, ok := err.(Error); ok {
		servErr.Write(ctx)
	} else {
		ctx.Logger().Error(err.Error())
		ErrInternalServer.Write(ctx)
	}
}

func (ctx *Ctx) nextHandler() Handler {
	ctx.indexHandler++
	if ctx.indexHandler < len(ctx.route.handlers) {
		return ctx.route.handlers[ctx.indexHandler]
	}
	return nil
}

func (ctx *Ctx) ParseInputJSON(dst any) error {
	return json.NewDecoder(ctx.Request().Body).Decode(dst)
}

func (ctx *Ctx) FormParam(field string) string {
	if !ctx.formParsed {
		err := ctx.Request().ParseMultipartForm(1024 * 1024)
		if err != nil {
			panic(err)
		}
		ctx.formParsed = true
	}
	return ctx.Request().FormValue(field)
}

func (ctx *Ctx) ContentType(contentType string) {
	ctx.SetHeader("Content-Type", contentType)
}

func (ctx *Ctx) SetHeader(name, value string) {
	if ctx.Sent() {
		ctx.Logger().Warn("response headers already sent", "header", name)
		return
	}
	ctx.w.Header().Set(name, value)
}

func (ctx *Ctx) Redirect(code int, url string) error {
	http.Redirect(ctx.w, ctx.r, url, code)
	return nil
}

func (ctx *Ctx) BadRequest(err error) error {
	if IsHTTPError(err) {
		return err
	}
	return NewError(http.StatusBadRequest, err.Error())
}

func (ctx *Ctx) NoContent() error {
	ctx.w.WriteHeader(http.StatusNoContent)
	return nil
}

func (ctx *Ctx) Bytes(code int, data []byte) error {
	ctx.w.WriteHeader(code)
	ctx.w.Write(data)
	return nil
}

func (ctx *Ctx) Copy(code int, r io.Reader) error {
	ctx.w.WriteHeader(code)
	_, err := io.Copy(ctx.w, r)
	return err
}

func (ctx *Ctx) Text(code int, text string) error {
	ctx.w.WriteHeader(code)
	io.WriteString(ctx.w, text)
	return nil
}

func (ctx *Ctx) JSON(code int, entity any) error {
	ctx.ContentType("application/json; charset=utf-8")
	ctx.w.WriteHeader(code)
	return json.NewEncoder(ctx.w).Encode(entity)
}

func (ctx *Ctx) ServeFile(file string) error {
	http.ServeFile(ctx.w, ctx.r, file)
	return nil
}
