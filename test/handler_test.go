package test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nskforward/httpx"
)

func TestHandler(t *testing.T) {
	app := httpx.New()

	app.Use(buildMiddleware("app-mw-1"), buildMiddleware("app-mw-2"))

	app.POST("/pass", func(c *httpx.Ctx) error {
		return c.Text(200, "pass")
	}, buildMiddleware("route-mw-1"), buildMiddleware("route-mw-2"))

	ts := httptest.NewServer(app.Handler())
	defer ts.Close()

	client := ts.Client()

	res, err := client.Post(ts.URL+"/pass", "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s: %s\n", res.Status, greeting)

	res, err = client.Get(ts.URL + "/404")
	if err != nil {
		t.Fatal(err)
	}
	greeting, err = io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s: %s\n", res.Status, greeting)

	res, err = client.Get(ts.URL + "/pass")
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

func buildMiddleware(text string) httpx.Handler {
	return func(c *httpx.Ctx) error {
		fmt.Println("b:", text)
		err := c.Next()
		fmt.Println("a:", text)
		return err
	}
}
