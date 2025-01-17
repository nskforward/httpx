package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
)

func TestMiddlewareBefore(t *testing.T) {
	r := httpx.NewRouter(Logger())

	r.Use(
		func(next httpx.Handler) httpx.Handler {
			return func(ctx *httpx.Context) error {
				ctx.Logger().Info("mw-1b")
				err := next(ctx)
				ctx.Logger().Info("mw-1a")
				return err
			}
		},
		func(next httpx.Handler) httpx.Handler {
			return func(ctx *httpx.Context) error {
				ctx.Logger().Info("mw-2b")
				err := next(ctx)
				ctx.Logger().Info("mw-2a")
				return err
			}
		},
	)

	r.Use(
		func(next httpx.Handler) httpx.Handler {
			return func(ctx *httpx.Context) error {
				ctx.Logger().Info("mw-3b")
				err := next(ctx)
				ctx.Logger().Info("mw-3a")
				return err
			}
		},
		func(next httpx.Handler) httpx.Handler {
			return func(ctx *httpx.Context) error {
				ctx.Logger().Info("mw-4b")
				err := next(ctx)
				ctx.Logger().Info("mw-4a")
				return err
			}
		},
	)

	r.GET("/",
		func(ctx *httpx.Context) error {
			ctx.Logger().Info("handler-1")
			return ctx.RespondText(200, "answer from server")
		},
		func(next httpx.Handler) httpx.Handler {
			return func(ctx *httpx.Context) error {
				ctx.Logger().Info("mw-5b")
				err := next(ctx)
				ctx.Logger().Info("mw-5a")
				return err
			}
		},
		func(next httpx.Handler) httpx.Handler {
			return func(ctx *httpx.Context) error {
				ctx.Logger().Info("mw-6b")
				err := next(ctx)
				ctx.Logger().Info("mw-6a")
				return err
			}
		},
	)

	s := httptest.NewServer(r)
	defer s.Close()

	res, err := http.Get(s.URL + "/test")
	if err != nil {
		t.Fatal(err)
	}

	answer, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(answer))
}
