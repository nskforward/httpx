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

func TestError(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println("--- main handler")
		return response.JSON(w, 200, response.H{"status": 100})
	}
	mw := func(msg any, after bool, hasError bool) types.Middleware {
		return func(next types.Handler) types.Handler {
			return func(w http.ResponseWriter, r *http.Request) error {
				if !after {
					fmt.Println("--- middleware:", msg)
				}
				if hasError {
					return response.Error{Status: http.StatusUnauthorized}
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

	r.Use(mw("1", false, false), mw("2", false, false))
	r.Use(mw("3", false, true))

	r.Route("/api/v1/", h, mw("4", false, false), mw("5", false, false))

	s := httptest.NewServer(&r)
	defer s.Close()

	DoRequest(s, "GET", "/api/v1/user/123", "", nil)
}
