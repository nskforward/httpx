package middleware

import (
	"net/http"

	"github.com/nskforward/httpx/cache"
	"github.com/nskforward/httpx/types"
)

func ETag(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		inm := r.Header.Get("If-None-Match")
		if inm != "" {
			etag, updated := cache.SearchETag(r)
			if etag == inm {
				w.Header().Set("Last-Modified", updated.UTC().Format(http.TimeFormat))
				w.Header().Set("ETag", inm)
				w.WriteHeader(http.StatusNotModified)
				return nil
			}
		}
		ims := r.Header.Get("If-Modified-Since")
		if ims != "" {
			_, updated := cache.SearchETag(r)
			if !updated.IsZero() && updated.UTC().Format(http.TimeFormat) == ims {
				w.Header().Set("Last-Modified", ims)
				w.WriteHeader(http.StatusNotModified)
				return nil
			}
		}

		return next(w, r)
	}
}
