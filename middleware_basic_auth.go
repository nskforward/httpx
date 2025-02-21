package httpx

func BasicAuth(usersWithPass map[string]string) Handler {
	return func(ctx *Ctx) error {
		user, pass, ok := ctx.Request().BasicAuth()
		if !ok || usersWithPass[user] != pass {
			ctx.SetHeader("WWW-Authenticate", `Basic realm="please provide credentials"`)
			return ErrUnauthorized
		}
		return ctx.Next()
	}
}
