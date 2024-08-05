package middleware

import (
	"net/http"

	"github.com/nskforward/httpx/cache"
	"github.com/nskforward/httpx/types"
)

func ETag(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "GET" && r.Method != "HEAD" {
			return next(w, r)
		}

		etag := r.Header.Get(types.IfNoneMatch)
		if etag != "" && etag == cache.LoadEtag(r) {
			cache.NotModifiend(w, etag)
		}

		return next(w, r)
	}
}

/*
{
	'Content-Type': 'text/event-stream',
	'Connection': 'keep-alive',
	'Cache-Control': 'no-cache',
	'X-Accel-Buffering': 'no'
}
*/
