package mux

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMux(t *testing.T) {
	mx := NewMultiplexer()
	mx.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "wildcard ANY /users/")
	})
	mx.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "exactly GET /")
	})
	mx.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "exactly GET /users")
	})
	mx.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("exactly DELETE /users id:%s", r.PathValue("id")))
	})
	mx.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("exactly GET /users id:%s", r.PathValue("id")))
	})
	mx.HandleFunc("GET /users/{id}/{action}", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("exactly GET /users id:%s action:%s", r.PathValue("id"), r.PathValue("action")))
	})

	// mx.Dump()

	r := httptest.NewRequest("GET", "/users/123/add", nil)
	t1 := time.Now()
	body, err := getResponseBody(mx, r)
	t2 := time.Since(t1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r.URL.Path, "-->", body, t2)
}

func getResponseBody(m *Multiplexer, r *http.Request) (string, error) {
	w := httptest.NewRecorder()
	h, code := m.Search(w, r)
	if code > 0 {
		return "", fmt.Errorf("bad response code: %d", code)
	}
	h.ServeHTTP(w, r)
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
