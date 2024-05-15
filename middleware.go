package httpx

type Middleware func(next Handler) Handler
