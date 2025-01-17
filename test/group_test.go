package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nskforward/httpx"
)

func TestGroup(t *testing.T) {
	r := httpx.NewRouter(Logger())

	users := r.Group("/users")
	users.Use(func(next httpx.Handler) httpx.Handler {
		return func(ctx *httpx.Context) error {
			ctx.Logger().Info("authorize user")
			return next(ctx)
		}
	})
	users.GET("/id/{id}", func(ctx *httpx.Context) error {
		return ctx.RespondText(200, fmt.Sprintf("user id: %s", ctx.PathParam("id")))
	})
	r.GET("/id/{id}", func(ctx *httpx.Context) error {
		return ctx.RespondText(200, fmt.Sprintf("guest id: %s", ctx.PathParam("id")))
	})

	answer := Exec(r, func(host string) (*http.Response, error) {
		return http.Get(host + "/users/id/123")
	})

	fmt.Println(string(answer))
}
