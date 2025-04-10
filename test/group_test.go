package test

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nskforward/httpx"
)

func TestGroup(t *testing.T) {
	app := httpx.NewRouter(slog.New(slog.NewTextHandler(io.Discard, nil)))

	app.Use(buildMiddleware("router-mw-1"), buildMiddleware("router-mw-2"))

	api := app.Group("/api", buildMiddleware("api-mw-1"), buildMiddleware("api-mw-2"))
	v1 := api.Group("/v1", buildMiddleware("v1-mw-1"), buildMiddleware("v1-mw-2"))

	v1.GET("/users", func(req *http.Request, resp *httpx.Response) error {
		return resp.Text(200, "GET /api/v1/users")
	})

	ts := httptest.NewServer(app)
	defer ts.Close()

	client := ts.Client()

	res, err := client.Get(ts.URL + "/api/v1/users")
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("client got: %d: %s\n", res.StatusCode, strings.TrimSpace(string(data)))
}
