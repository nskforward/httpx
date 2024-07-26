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

func DoRequest(s *httptest.Server, method, uri, body string, header http.Header, dumpHeaders, dumpBody bool) {
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

	req.AddCookie(&http.Cookie{
		Name:  "SID",
		Value: "123",
	})

	t1 := time.Now()
	resp, err := http.DefaultClient.Do(req)
	t2 := time.Since(t1)
	if err != nil {
		panic(err)
	}
	if dumpHeaders {
		data, err := httputil.DumpResponse(resp, dumpBody)
		if err != nil {
			panic(err)
		}
		fmt.Println()
		fmt.Println(t2)
		fmt.Println(string(data))
	}
}

func doReq(method, urlString string, header http.Header) *http.Response {
	req, err := http.NewRequest(method, urlString, nil)
	if err != nil {
		panic(err)
	}
	for k, vv := range header {
		for _, v := range vv {
			req.Header.Set(k, v)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	data, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println(string(data))
	return resp
}
