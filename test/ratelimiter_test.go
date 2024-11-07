package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nskforward/httpx"
	m "github.com/nskforward/httpx/middleware"
	"github.com/nskforward/httpx/proxy"
	"github.com/nskforward/httpx/types"
)

func TestRateLimiter(t *testing.T) {

	// BACKEND
	backend1 := httptest.NewServer(
		httpx.NewRouter().Route("/", httpx.Text("success"), m.RealIP, m.SetHeader("Server", "backend", false)),
	)
	defer backend1.Close()

	// PROXY
	r := httpx.NewRouter()
	r.Use(m.SetHeader("Server", "proxy", true), m.RateLimiter(
		time.Second,
		1,
		2*time.Second,
		func(r *http.Request) string {
			return r.URL.Path
		}),
	)
	r.Route("/", proxy.Reverse(backend1.URL))

	frontendProxy := httptest.NewServer(r)
	defer frontendProxy.Close()

	for range 5 {
		go DoRequest(frontendProxy, "POST", "/test", "aaa=111", http.Header{types.AcceptEncoding: []string{"gzip"}}, false, false)
	}

	time.Sleep(3 * time.Second)
}
