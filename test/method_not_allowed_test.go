package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nskforward/httpx"
)

func TestMethodNotAllowed(t *testing.T) {
	router := httpx.NewRouter(nil)
	router.Use(func(req *http.Request, resp *httpx.Response) error {
		fmt.Println("event:", time.Now().Format("15:04:05.000"), req.Method, req.URL.Path)
		err := resp.Next()
		fmt.Println("status:", resp.StatusCode())
		if err != nil {
			fmt.Println("error:", err)
		}
		return nil
	})
	router.CustomMethod("POST", "/get", func(req *http.Request, resp *httpx.Response) error {
		return resp.Text(200, "hello from handler")
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()

	res, err := client.Get(ts.URL + "/get")
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("client got result: status (%s) body (%s)\n", res.Status, strings.TrimSpace(string(data)))
}
