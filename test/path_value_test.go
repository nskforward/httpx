package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPathValue(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("id:", r.PathValue("id"))
		w.WriteHeader(http.StatusNoContent)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := ts.Client()

	req, err := http.NewRequest("DELETE", ts.URL+"/api/v1/user/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode/100 != 2 {
		t.Fatal(fmt.Errorf("bad response code: %s", resp.Status))
	}
}
