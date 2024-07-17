package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/response"
)

func TestParams(t *testing.T) {

	var r httpx.Router

	r.Route("/api/{version}/user/{id}/{path...}", func(w http.ResponseWriter, r *http.Request) error {
		return response.Text(w, 200, fmt.Sprintf("%s, %s, %s", r.PathValue("version"), r.PathValue("id"), r.PathValue("path")))
	})

	s := httptest.NewServer(&r)
	defer s.Close()

	DoRequest(s, "GET", "/api/v1/user/123/filter/admin", "")
}
