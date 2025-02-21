package httpx

func HealthCheck(method, path string) Handler {
	return func(ctx *Ctx) error {
		if ctx.Request().Method == method && ctx.Path() == path {
			return ctx.Text(200, "ok")
		}
		return ctx.Next()
	}
}
