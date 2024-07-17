package middleware

import (
	"net/http"
	"strings"

	"github.com/nskforward/httpx/types"
)

func SanitizePath(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if len(r.URL.Path) == 0 {
			r.URL.Path = "/"
			r.URL.RawPath = "/"
			return next(w, r)
		}
		if r.URL.Path[0] != '/' {
			r.URL.Path = "/" + r.URL.Path
			r.URL.RawPath = "/" + r.URL.RawPath
		}
		if len(r.URL.Path) > 1 && r.URL.Path[len(r.URL.Path)-1] == '/' {
			r.URL.Path = strings.TrimRight(r.URL.Path, "/")
			r.URL.RawPath = strings.TrimRight(r.URL.RawPath, "/")
		}
		return next(w, r)
	}
}
