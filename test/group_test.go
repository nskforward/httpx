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

func TestGroup(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println("--- main handler")
		return response.Text(w, 200, "success")
	}
	mw := func(msg any, after bool) types.Middleware {
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

	r.Use(mw(0, false))

	group1 := r.Group(mw(1, false))
	group1.Route("/api/v1/", h, mw(11, false), mw(12, false))

	group2 := r.Group(mw(2, false))
	group2.Route("/api/v2/", h, mw(21, false), mw(22, false))

	r.Route("/api/", h, mw(3, false))

	s := httptest.NewServer(&r)
	defer s.Close()

	DoRequest(s, "GET", "/api/v1/user/111", "", nil)
	DoRequest(s, "GET", "/api/v2/user/222", "", nil)
	DoRequest(s, "GET", "/api/v3/useår/333", "", nil)
}
