package test

import (
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/middleware"
)

func TestFileDir(t *testing.T) {
	var r httpx.Router
	r.Use(middleware.RequestID)
	r.Route("/", httpx.ServeDir("data"))

	s := httptest.NewServer(&r)
	defer s.Close()

	DoRequest(s, "GET", "/static", "", nil)
}
