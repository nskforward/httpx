package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"time"
)

func HTTPClient() *http.Client {
	return &http.Client{Transport: &http.Transport{DisableCompression: true}}
}

func DoRequest(s *httptest.Server, method, uri, body string, header http.Header) {
	fmt.Println("================================================")
	fmt.Println(method, uri)
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, s.URL+uri, reader)
	if err != nil {
		panic(err)
	}
	for k, vv := range header {
		for _, v := range vv {
			req.Header.Set(k, v)
		}
	}
	t1 := time.Now()
	resp, err := http.DefaultClient.Do(req)
	t2 := time.Since(t1)
	if err != nil {
		panic(err)
	}
	data, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println(t2)
	fmt.Println(string(data))
}
