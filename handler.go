package httpx

import "net/http"

type Handler func(req *http.Request, resp *Response) error
