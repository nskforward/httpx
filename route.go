package httpx

type Route struct {
	pattern     string
	handler     Handler
	middlewares []Middleware
}
