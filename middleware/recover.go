package middleware

import (
	"fmt"
	"net/http"

	"github.com/nskforward/httpx"
)

func Recover(req *http.Request, resp *httpx.Response) error {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			if err == http.ErrAbortHandler {
				panic(err)
			}
			resp.ServerError(err)
		}
	}()

	return resp.Next(req)
}
