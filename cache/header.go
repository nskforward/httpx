package cache

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nskforward/httpx/types"
)

const MaxAge = 31536000 * time.Second

func HeaderNoCache(w http.ResponseWriter, etag string, updated time.Time) {
	w.Header().Set(types.CacheControl, "no-cache")
	w.Header().Set(types.LastModified, updated.UTC().Format(http.TimeFormat))
	w.Header().Set(types.ETag, etag)
}

func HeaderNoStoreCache(w http.ResponseWriter) {
	w.Header().Set(types.CacheControl, "no-store")
}

func HeaderCache(w http.ResponseWriter, shared bool, age time.Duration) {
	w.Header().Set(types.LastModified, time.Now().UTC().Format(http.TimeFormat))

	if shared {
		w.Header().Set(types.CacheControl, "public")
	} else {
		w.Header().Set(types.CacheControl, "private")
	}

	if age > MaxAge {
		age = MaxAge
	}

	w.Header().Set(types.CacheControl, fmt.Sprintf("max-age=%.0f", age.Seconds()))

	if age == MaxAge {
		w.Header().Set(types.CacheControl, "immutable")
	}
}
