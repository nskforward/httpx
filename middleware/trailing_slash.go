package middleware

import (
	"net/http"
	"strings"

	"github.com/nskforward/httpx/types"
)

type TrailingSlashAction uint8

const (
	Redirect TrailingSlashAction = iota
	Continue
)

func TrailingSlash(action TrailingSlashAction) types.Middleware {
	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			has := len(r.URL.Path) > 1 && r.URL.Path[len(r.URL.Path)-1] == '/'
			if has {
				if action == Redirect {
					http.Redirect(w, r, strings.TrimRight(r.URL.RequestURI(), "/"), http.StatusPermanentRedirect)
					return nil
				} else {
					r.URL.Path = strings.TrimRight(r.URL.Path, "/")
					r.URL.RawPath = strings.TrimRight(r.URL.RawPath, "/")
				}
			}
			return next(w, r)
		}
	}
}
