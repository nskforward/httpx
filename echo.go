package httpx

import (
	"fmt"
	"io"
	"net/http"
)

func Echo(w http.ResponseWriter, r *http.Request) error {
	io.WriteString(w, fmt.Sprintf("%s %s %s\n", r.Method, r.URL.RequestURI(), r.Proto))
	for k, vv := range r.Header {
		for _, v := range vv {
			io.WriteString(w, fmt.Sprintf("%s: %s\n", k, v))
		}
	}
	io.WriteString(w, "\n")
	io.Copy(w, r.Body)
	return nil
}
