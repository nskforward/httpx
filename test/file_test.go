package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/middleware"
	"github.com/nskforward/httpx/types"
)

func TestFile(t *testing.T) {
	var r httpx.Router
	r.Use(middleware.RequestID)
	r.Route("/api/v1/", httpx.ServeFile("data/static/123.html"))

	s := httptest.NewServer(&r)
	defer s.Close()

	DoRequest(s, "GET", "/api/v1/user/123", "", http.Header{types.AcceptEncoding: []string{"gzip"}}, true, false)
}
