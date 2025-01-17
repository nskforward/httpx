package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func Exec(r http.Handler, request func(serverHost string) (*http.Response, error)) []byte {
	s := httptest.NewServer(r)
	defer s.Close()

	resp, err := request(s.URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	answer, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode/100 != 2 {
		panic(fmt.Errorf("bad response status: %s", resp.Status))
	}

	return answer
}
