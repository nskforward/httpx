package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/middleware"
	"github.com/nskforward/httpx/proxy"
	"github.com/nskforward/httpx/types"
)

func TestProxySingle(t *testing.T) {

	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "backend-1")
	}))
	defer backend1.Close()

	var r httpx.Router
	r.Use(middleware.Log, middleware.Recovery)
	r.Route("/api/v1/", proxy.Reverse(backend1.URL))

	frontendProxy := httptest.NewServer(&r)
	defer frontendProxy.Close()

	DoRequest(frontendProxy, "GET", "/api/v1/user/123", "", http.Header{types.AcceptEncoding: []string{"gzip"}})
}
