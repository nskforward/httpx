package httpx

type Router struct {
	multiplexer Multiplexer
}

func NewRouter(multiplexer Multiplexer) *Router {
	return &Router{
		multiplexer: multiplexer,
	}
}
