package middleware

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nskforward/httpx/cache"
	"github.com/nskforward/httpx/types"
)

func Cache(dir string, maxSize cache.Size) types.Middleware {

	settings, err := cache.ValidateSettings(cache.Settings{
		Dir:          dir,
		TotalMaxSize: maxSize,
	})
	if err != nil {
		panic(err)
	}

	store := cache.NewStore(settings.Dir)

	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			// in any handlers we can get this store and set the cache
			r = store.Inject(r)

			if r.Method != http.MethodGet {
				return next(w, r)
			}

			if handleCachedResponse(store, w, r) {
				// hit cache
				return nil
			}

			// miss cache
			var f *os.File
			var entry *cache.Entry
			ww := types.NewResponseWrapper(w)
			ww.BeforeBody = func() {
				if ww.Status() != 200 {
					return
				}

				ttl := getTTL(r, ww.Header())
				if ttl < 0 {
					return
				}

				ww.Header().Set(types.XCache, "miss")

				bucket := store.GetOrCreateBucket(r.URL.RequestURI())
				entry = bucket.GetKey(r)
				if entry == nil {
					entry = bucket.SetKey(r, ttl, func(r *http.Request) string {
						return r.URL.RequestURI()
					})
				}
				f, err = os.Create(entry.File())
				if err != nil {
					slog.Error("cannot create a cache file", "error", err)
					return
				}
				if !entry.SetFilling() {
					return
				}
				writeDownHeaders(f, ww.Status(), ww.Header())
				ww.SetWriter(io.MultiWriter(ww.ResponseWriter, f))
			}
			err := next(ww, r)
			if entry != nil {
				entry.SetIdle()
			}
			if f != nil {
				f.Close()
			}
			return err
		}
	}
}

func writeDownHeaders(f *os.File, status int, header http.Header) {
	io.WriteString(f, strconv.Itoa(status))
	io.WriteString(f, "\n")
	for k, vv := range header {
		if k == "Connection" {
			continue
		}
		if k == "Accept-Ranges" {
			continue
		}
		io.WriteString(f, k)
		io.WriteString(f, ": ")
		io.WriteString(f, strings.Join(vv, ", "))
		io.WriteString(f, "\n")
	}
	io.WriteString(f, "\n")
}

func getTTL(r *http.Request, respheader http.Header) time.Duration {
	cacheControl := cache.NewCacheControl(respheader.Get(types.CacheControl))

	if cacheControl.NoStore && cacheControl.Private {
		return -1
	}
	if r.Header.Get(types.Authorization) != "" && !cacheControl.Public {
		return -1
	}
	if cacheControl.SMaxAge > 0 {
		return cacheControl.SMaxAge * time.Second
	}
	if cacheControl.MaxAge > 0 {
		return cacheControl.MaxAge * time.Second
	}
	if cacheControl.NoCache {
		return 0
	}
	expires, err := time.Parse(http.TimeFormat, respheader.Get(types.Expires))
	if err == nil && !expires.IsZero() {
		return time.Until(expires)
	}
	lastModified, err := time.Parse(http.TimeFormat, respheader.Get(types.LastModified))
	if err == nil && !lastModified.IsZero() && time.Since(lastModified) > 10*time.Second {
		return time.Since(lastModified) / 10
	}
	return -1
}

func handleCachedResponse(store *cache.Store, w http.ResponseWriter, r *http.Request) bool {
	bucket := store.GetBucket(r.URL.RequestURI())
	if bucket == nil {
		return false
	}
	entry := bucket.GetKey(r)
	if entry == nil {
		return false
	}
	if !entry.IsIdle() {
		return false
	}

	etag := r.Header.Get(types.IfNoneMatch)
	if etag != "" && etag == entry.ID() {
		w.Header().Set(types.ETag, etag)
		w.Header().Set(types.LastModified, entry.LastModified())
		w.WriteHeader(http.StatusNotModified)
		return true
	}

	if !entry.Valid() {
		return false
	}

	err := sendCache(w, entry)
	if err != nil {
		slog.Error("sending cache file", "error", err)
		return false
	}

	return true
}

func sendCache(w http.ResponseWriter, entry *cache.Entry) error {
	f, err := os.Open(entry.File())
	if err != nil {
		return err
	}
	defer f.Close()
	b := bufio.NewReader(f)
	line, err := b.ReadBytes('\n')
	if err != nil {
		return err
	}
	line = bytes.TrimRight(line, "\n")
	status, err := strconv.Atoi(string(line))
	if err != nil {
		return err
	}
	if status < 100 || status > 599 {
		return fmt.Errorf("bad status format: %s", string(line))
	}
	for len(line) > 0 {
		line, err := b.ReadBytes('\n')
		if err != nil {
			return err
		}
		line = bytes.TrimRight(line, "\n")
		header := bytes.Split(line, []byte(": "))
		if len(header) != 2 {
			return fmt.Errorf("bad header format: %s", string(line))
		}
		w.Header().Set(string(line[0]), string(line[1]))
	}
	w.Header().Set("X-Cache", "hit")
	w.Header().Set("Accept-Ranges", "none")
	w.Header().Set("Age", strconv.FormatFloat(time.Since(entry.From()).Seconds(), 'f', 0, 64))
	io.Copy(w, b)
	return nil
}
