package httpx

// Handler returns an error only in unexpected behaviors (this error will be sent as internal server error with status 500)
type Handler func(ctx *Context) error

type Middleware func(next Handler) Handler

type LogFunc func(ctx *Context)
