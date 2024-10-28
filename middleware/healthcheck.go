package middleware

import (
	"net/http"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func Healthcheck(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodGet && r.URL.Path == "/healthcheck" {
			return response.Text(w, 200, "ok")
		}
		return next(w, r)
	}
}
