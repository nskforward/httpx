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
	var r1 httpx.Router
	r1.Route("/api/v1/", httpx.Echo, middleware.RealIP)

	backend1 := httptest.NewServer(&r1)
	defer backend1.Close()

	var r2 httpx.Router
	r2.Route("/api/v1/", proxy.Reverse(backend1.URL))

	frontendProxy := httptest.NewServer(&r2)
	defer frontendProxy.Close()

	fmt.Println("proxy:", frontendProxy.URL)
	fmt.Println("backend:", backend1.URL)

	DoRequest(frontendProxy, "POST", "/api/v1/user/123", "aaa=111", http.Header{types.AcceptEncoding: []string{"gzip"}})
}
