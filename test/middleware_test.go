package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nskforward/httpx"
)

func TestMiddleware(t *testing.T) {
	router := httpx.NewRouter(nil)
	router.Use(buildMiddleware("router-mw-1"))

	router.GET("/", func(req *http.Request, resp *httpx.Response) error {
		fmt.Println("call handler")
		return resp.Text(200, "pass")
	}, buildMiddleware("handler-mw-1"))

	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()
	response, err := client.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("client got: %d: %s\n", response.StatusCode, strings.TrimSpace(string(data)))
}

func buildMiddleware(text string) httpx.Handler {
	return func(req *http.Request, resp *httpx.Response) error {
		fmt.Println("call middleware before:", text)
		err := resp.Next()
		fmt.Println("call middleware after:", text)
		return err
	}
}
