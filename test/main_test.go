package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/response"
)

func TestMain(t *testing.T) {

	testCases := []struct {
		pattern    string
		answer     string
		status     int
		reqMethod  string
		reqURI     string
		respStatus int
		respText   string
	}{
		{
			pattern:    "GET /api/v1/user1",
			status:     200,
			reqMethod:  "GET",
			reqURI:     "/api/v1/user1",
			respStatus: 200,
			respText:   "1:/api/v1/user1",
		},
		{
			pattern:    "POST /api/v1/user2",
			status:     200,
			reqMethod:  "GET",
			reqURI:     "/api/v1/user2",
			respStatus: 200,
			respText:   "5:/api/v1/user2",
		},
		{
			pattern:    "/api/v1/foo/",
			status:     200,
			reqMethod:  "GET",
			reqURI:     "/api/v1/foo/user3",
			respStatus: 200,
			respText:   "3:/api/v1/foo/user3",
		},
		{
			pattern:    "/tmp/1",
			status:     200,
			reqMethod:  "GET",
			reqURI:     "/api/v1/foo/user4",
			respStatus: 200,
			respText:   "3:/api/v1/foo/user4",
		},
		{
			pattern:    "/",
			status:     200,
			reqMethod:  "GET",
			reqURI:     "/tmp/2",
			respStatus: 200,
			respText:   "5:/tmp/2",
		},
		{
			pattern:    "/{$}",
			status:     200,
			reqMethod:  "GET",
			reqURI:     "/",
			respStatus: 200,
			respText:   "6:/",
		},
	}

	r := httpx.NewRouter()
	for i, tc := range testCases {
		r.Route(tc.pattern, func(w http.ResponseWriter, r *http.Request) error {
			return response.Text(w, tc.status, fmt.Sprintf("%d:%s", i+1, r.URL.RequestURI()))
		})
	}

	s := httptest.NewServer(r)
	defer s.Close()

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.reqMethod, s.URL+tc.reqURI, nil)
		if err != nil {
			t.Fatal(err)
		}
		t1 := time.Now()
		resp, err := http.DefaultClient.Do(req)
		t2 := time.Since(t1)
		if err != nil {
			t.Fatal(err)
		}
		data, err := httputil.DumpResponse(resp, false)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("================================================")
		fmt.Println(t2)
		fmt.Print(string(data))

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		fmt.Println(string(body))

		if resp.StatusCode != tc.respStatus {
			t.Fatalf("URI '%s' expect status code %d, actual %d", tc.reqURI, tc.respStatus, resp.StatusCode)
		}

		if string(body) != tc.respText {
			t.Fatalf("URI '%s' expect body '%s', actual '%s'", tc.reqURI, tc.respText, string(body))
		}
	}
}
