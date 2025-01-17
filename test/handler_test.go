package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nskforward/httpx"
)

func TestHandler(t *testing.T) {
	r := httpx.NewRouter(Logger())

	r.GET("/test", func(ctx *httpx.Context) error {
		return ctx.RespondText(200, "get!")
	})

	r.DELETE("/test", func(ctx *httpx.Context) error {
		return ctx.RespondText(200, "delete!")
	})

	s := httptest.NewServer(r)
	defer s.Close()

	res, err := http.Get(s.URL + "/test")
	if err != nil {
		t.Fatal(err)
	}

	answer, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(answer))
}
