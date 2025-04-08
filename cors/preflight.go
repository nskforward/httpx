package cors

import "net/http"

func isPreflight(req *http.Request) bool {
	if req.Method != http.MethodOptions || req.Header.Get("Origin") == "" {
		return false
	}
	return req.Header.Get("Access-Control-Request-Method") != "" || req.Header.Get("Access-Control-Request-Headers") != ""
}
