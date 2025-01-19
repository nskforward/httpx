package httpx

func BasicAuthMiddleware(users map[string]string) Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) error {
			user, pass, ok := ctx.Request().BasicAuth()
			if !ok || users[user] != pass {
				ctx.SetResponseHeader("WWW-Authenticate", `Basic realm="please provide app credentials"`)
				return ctx.Unauthorized()
			}
			return next(ctx)
		}
	}
}
