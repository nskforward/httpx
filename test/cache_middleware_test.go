package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/cache"
	m "github.com/nskforward/httpx/middleware"
	"github.com/nskforward/httpx/proxy"
	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func TestCacheMiddleware(t *testing.T) {
	// BACKEND
	br := httpx.NewRouter()

	br.Use(m.RealIP, m.SetHeader("Server", "backend", false), m.Cache("data/cache", cache.GB))

	br.Route("/api/v1/users", func(w http.ResponseWriter, r *http.Request) error {
		store := cache.GetStore(r)
		if store != nil {
			bucket := store.GetOrCreateBucket(r.URL.RequestURI())
			bucket.Tag("users")
			entry := bucket.SetKey(r, 0, func(r *http.Request) string {
				return r.URL.RequestURI()
			})
			entry.SendNoCache(w)
		}
		return response.JSON(w, 200, []string{"ivan", "alex", "oleg", "john"})
	})

	br.Route("POST /api/v1/users", func(w http.ResponseWriter, r *http.Request) error {
		store := cache.GetStore(r)
		if store != nil {
			tag := store.GetTag("users")
			tag.Delete()
		}
		return response.NoContent(w)
	})
	backend := httptest.NewServer(br)
	defer backend.Close()

	// PROXY
	pr := httpx.NewRouter()
	pr.Use(m.Recover, m.RequestID)
	pr.Route("/api/v1/", proxy.Reverse(backend.URL), m.SetHeader("Server", "proxy", true))
	frontend := httptest.NewServer(pr)
	defer frontend.Close()

	resp := doReq("GET", frontend.URL+"/api/v1/users", nil)
	reqHeader := http.Header{}
	if resp.Header.Get(types.ETag) != "" {
		reqHeader.Set(types.IfNoneMatch, resp.Header.Get(types.ETag))
	}
	DoRequest(frontend, "GET", "/api/v1/users", "", reqHeader, true, true)
	time.Sleep(time.Second)
	DoRequest(frontend, "POST", "/api/v1/users", "", nil, true, true)
	time.Sleep(time.Second)
	DoRequest(frontend, "GET", "/api/v1/users", "", reqHeader, true, true)
}
