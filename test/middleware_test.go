package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func TestMiddleware(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println("--- main handler")
		return response.Text(w, 200, "success")
	}
	mw := func(msg string, after bool) types.Middleware {
		return func(next types.Handler) types.Handler {
			return func(w http.ResponseWriter, r *http.Request) error {
				if !after {
					fmt.Println("--- middleware:", msg)
				}
				err := next(w, r)
				if after {
					fmt.Println("--- middleware:", msg)
				}
				return err
			}
		}
	}

	var r httpx.Router

	r.Use(mw("1", false), mw("2", false))
	r.Use(mw("3", false))

	r.Route("/api/v1/", h, mw("4", false), mw("5", false))

	s := httptest.NewServer(&r)
	defer s.Close()

	DoRequest(s, "GET", "/api/v1/user/123", "")
}
