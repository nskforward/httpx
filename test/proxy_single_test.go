package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
	m "github.com/nskforward/httpx/middleware"
	"github.com/nskforward/httpx/proxy"
	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func TestProxySingle(t *testing.T) {

	// BACKEND
	br := httpx.NewRouter()
	br.Use(m.RealIP, m.SetHeader("Server", "backend", false))
	br.Route("/", httpx.Echo)
	br.Route("/cookies", func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println("cookies:", r.Header.Get("Cookie"))
		return response.Text(w, 200, "dump cookies")
	})
	backend := httptest.NewServer(br)
	defer backend.Close()

	// PROXY
	pr := httpx.NewRouter()
	pr.Use(m.Recover, m.TraceIDSetter)
	pr.Route("/", proxy.Reverse(backend.URL), m.SetHeader("Server", "proxy", true))
	pr.Route("/test", httpx.Text("hello from proxy!"), m.SetHeader("Server", "proxy", true))
	frontend := httptest.NewServer(pr)
	defer frontend.Close()

	DoRequest(frontend, "POST", "/cookies", "aaa=111", http.Header{types.AcceptEncoding: []string{"gzip"}}, true, true)
}
