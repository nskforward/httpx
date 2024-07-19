package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/middleware"
	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func TestRecover(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println("--- main handler")
		panic("my test panic")
		// return response.Error{Status: 500, Text: "main handler error triggered"}
		//return response.JSON(w, 200, response.H{"status": 100})
	}
	mw := func(msg any, after bool, hasError bool) types.Middleware {
		return func(next types.Handler) types.Handler {
			return func(w http.ResponseWriter, r *http.Request) error {
				if !after {
					fmt.Println("--- middleware:", msg)
				}
				if hasError {
					return response.Error{Status: http.StatusUnauthorized, Text: "test error"}
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

	r.Use(middleware.Recover, middleware.RequestID)
	r.Use(mw(1, false, false))

	r.Route("/api/v1/", h, mw(2, false, false))

	s := httptest.NewServer(&r)
	defer s.Close()

	DoRequest(s, "GET", "/api/v1/user/123", "", nil)
}
