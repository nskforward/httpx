package httpx

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestRouter(t *testing.T) {
	ro := NewRouter(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	ro.Use(getMiddleware("middleware 1"))

	group := ro.Group(getMiddleware("middleware 2"))

	ro.HandleFunc("/users/", getHandler("ANY /users/", nil), getMiddleware("middleware 3.1"))
	group.HandleFunc("POST /users", getHandler("POST /users", nil), getMiddleware("middleware 3.2"))
	ro.HandleFunc("GET /users/{action}/{id}", getHandler("GET /users/{action}/{id}", []string{"action", "id"}), getMiddleware("middleware 3.3"))

	r := httptest.NewRequest("POST", "/users", nil)
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
