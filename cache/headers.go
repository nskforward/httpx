package cache

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nskforward/httpx/types"
)

func NotModifiend(w http.ResponseWriter, etag string) {
	w.Header().Set(types.ETag, etag)
	w.WriteHeader(http.StatusNotModified)
}

func Prohibit(w http.ResponseWriter) {
	w.Header().Set(types.CacheControl, "no-store")
}

func Conditional(w http.ResponseWriter) {
	w.Header().Set(types.CacheControl, "no-cache")
}

func Public(w http.ResponseWriter, maxAge time.Duration) {
	w.Header().Set(types.CacheControl, fmt.Sprintf("public, max-age=%.0f", maxAge.Seconds()))
}

func Private(w http.ResponseWriter, maxAge time.Duration) {
	w.Header().Set(types.CacheControl, fmt.Sprintf("private, max-age=%.0f", maxAge.Seconds()))
}

func ETag(w http.ResponseWriter, r *http.Request, keyFunc KeyFunc) {
	etag := SaveEtag(r, keyFunc)
	if etag != "" {
		Conditional(w)
		w.Header().Set(types.ETag, etag)
	}
}
