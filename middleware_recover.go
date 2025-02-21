package httpx

import (
	"fmt"
	"net/http"
)

func Recover(ctx *Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)

			}
			if err == http.ErrAbortHandler {
				panic(err)
			}
			ctx.Logger().Error("panic", "error", err)
			ErrInternalServer.Write(ctx)
		}
	}()

	return ctx.Next()
}
