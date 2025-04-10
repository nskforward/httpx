package test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/nskforward/httpx/cors"
)

func TestCorsOrigin(t *testing.T) {
	origins := []struct {
		cors    string
		request string
		valid   bool
	}{
		{
			cors:    "http://example.com",
			request: "http://example.com",
			valid:   true,
		},
		{
			cors:    "https://sub.example.com",
			request: "https://sub.example.com",
			valid:   true,
		},
		{
			cors:    "https://example.com",
			request: "https://sub.example.com",
			valid:   false,
		},
		{
			cors:    "https://example.com",
			request: "http://example.com",
			valid:   false,
		},
		{
			cors:    "https://example.com:443",
			request: "https://example.com",
			valid:   true,
		},
		{
			cors:    "https://example.com",
			request: "https://example.com:443",
			valid:   true,
		},
		{
			cors:    "https://example.com:8081",
			request: "https://example.com",
			valid:   false,
		},
		{
			cors:    "https://example.com:8081",
			request: "https://example.com:8080",
			valid:   false,
		},
		{
			cors:    "https://*.example.com",
			request: "https://sub.example.com",
			valid:   true,
		},
		{
			cors:    "https://*.example.com",
			request: "https://example.com",
			valid:   false,
		},
		{
			cors:    "https://*.example.com",
			request: "https://sub.example2.com",
			valid:   false,
		},
		{
			cors:    "https://*.example.com",
			request: "https://sub.example.com:443",
			valid:   true,
		},
	}

	for index, item := range origins {
		origin, err := cors.ParseOrigin(item.cors)
		if err != nil {
			t.Fatal(fmt.Errorf("%d: fail: %w", index+1, err))
		}
		input, err := url.Parse(item.request)
		if err != nil {
			t.Fatal(fmt.Errorf("%d: fail: %w", index+1, err))
		}
		if origin.Valid(input) != item.valid {
			t.Fatal(fmt.Errorf("%d: fail: expect valid=%v", index+1, item.valid))
		}
		fmt.Printf("%d: success\n", index+1)
	}
}
