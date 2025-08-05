package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type contextParam string

const (
	userKey contextParam = "user"
)

func TestRouter(t *testing.T) {
	ro := NewRouter()
	ro.Use(func(next HandlerFunc) HandlerFunc {
		return func(w *Response, r *http.Request) error {
			return next(w, r.WithContext(context.WithValue(r.Context(), userKey, "Ivan")))
		}
	})

	ro.HandleFunc("/users/", func(w *Response, r *http.Request) error {
		user := r.Context().Value(userKey)
		if user == nil {
			w.SendError(400, "Bad Request")
			return fmt.Errorf("user cannot be empty")
		}
		io.WriteString(w, fmt.Sprintf("user: %s", user))
		return nil
	})

	r := httptest.NewRequest("GET", "/users/123", nil)
	t1 := time.Now()
	body, err := getResponseBody(ro, r)
	t2 := time.Since(t1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r.URL.Path, "-->", body, t2)
}

func getHandler(text string, params []string) HandlerFunc {
	return func(w *Response, r *http.Request) error {
		fmt.Println("call handler", text)
		for _, key := range params {
			val := r.PathValue(key)
			fmt.Println("Param", key, "=", val)
		}
		io.WriteString(w, text)
		return nil
	}
}

func getMiddleware(text string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(w *Response, r *http.Request) error {
			fmt.Println(text, "before")
			err := next(w, r)
			fmt.Println(text, "after")
			return err
		}
	}
}

func getResponseBody(ro *Router, r *http.Request) (string, error) {
	w := httptest.NewRecorder()
	ro.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode/100 != 2 {
		return "", fmt.Errorf("bad response code %d: %s", res.StatusCode, string(data))
	}
	return string(data), nil
}
