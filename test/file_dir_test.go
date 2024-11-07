package test

import (
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
)

func TestFileDir(t *testing.T) {
	var r httpx.Router
	r.Route("/", httpx.ServeDir("data"))

	s := httptest.NewServer(&r)
	defer s.Close()

	DoRequest(s, "GET", "/static", "", nil, true, false)
}
