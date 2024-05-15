package httpx

import (
	"net/http"
	"time"

	"github.com/nskforward/httpx/cache"
	"github.com/nskforward/httpx/cache/store"
)

const maxCacheFileSize = 1024 * 1024 // 1 MB

func Cache(cacheDir string, defaultTTL time.Duration) Middleware {

	s := store.NewStore(cacheDir, defaultTTL)

	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if r.Method != http.MethodGet {
				// allow caching the GET requests only
				return next(w, r)
			}

			key := r.URL.Path
			entry, err := s.Get(key)
			if err != nil {
				return InternalServerError(err)
			}
			if entry != nil {
				if !entry.Expired() {
					w.Header().Set("X-Cache", "HIT")
					entry.ServeHTTP(w, r)
					return nil
				}
			}

			ww := cache.GetWriter()
			ww.Reset(w, maxCacheFileSize)
			defer cache.PutWriter(ww)

			ww.Header().Set("X-Cache", "MISS")

			err = next(ww, r)
			if err != nil {
				return err
			}

			if !ww.CanCache() {
				return nil
			}

			ttl := defaultTTL
			if ww.CacheAge() > 0 {
				ttl = ww.CacheAge()
			}

			err = s.Set(key, ttl, ww.Buffer(), map[string]string{
				"Content-Type":     ww.Header().Get("Content-Type"),
				"Content-Encoding": ww.Header().Get("Content-Encoding"),
			})
			if err != nil {
				return InternalServerError(err)
			}

			return nil
		}
	}
}
