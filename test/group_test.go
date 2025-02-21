package test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nskforward/httpx"
)

func TestGroup(t *testing.T) {
	app := httpx.New()
	app.Use(buildMiddleware("app-mw-1"))
	app.Use(buildMiddleware("app-mw-2"))

	api := app.Group("/api", buildMiddleware("api-mw-1"), buildMiddleware("api-mw-2"))
	v1 := api.Group("/v1", buildMiddleware("v1-mw-1"), buildMiddleware("v1-mw-2"))

	v1.Route(httpx.GET, "/users", func(c *httpx.Ctx) error {
		return c.Text(200, "GET /api/v1/users")
	}, buildMiddleware("get-users-mw-1"), buildMiddleware("get-users-mw-2"))

	v1.Route(httpx.POST, "/users", func(c *httpx.Ctx) error {
		return c.Text(200, "POST /api/v1/users")
	}, buildMiddleware("post-users-mw-1"), buildMiddleware("post-users-mw-2"))

	ts := httptest.NewServer(app.Handler())
	defer ts.Close()

	client := ts.Client()

	res, err := client.Get(ts.URL + "/api/v1/users")
	if err != nil {
		t.Fatal(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s: %s\n", res.Status, greeting)

	res, err = client.Post(ts.URL+"/api/v1/users", "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	greeting, err = io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s: %s\n", res.Status, greeting)
}
