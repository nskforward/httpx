package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/middleware"
	"github.com/nskforward/httpx/response"
)

func TestSetHeader(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println("--- main handler")
		return response.JSON(w, 200, response.H{"status": "success"})
	}

	var r httpx.Router
	r.Use(middleware.Log, middleware.Recovery)
	r.Use(
		middleware.SetHeader("Server", "unit-test1", false),
		middleware.SetHeader("Server", "unit-test2", true),
	)
	r.Route("/api/v1/", h)

	s := httptest.NewServer(&r)
	defer s.Close()

	DoRequest(s, "GET", "/api/v1/user/123", "")
}
