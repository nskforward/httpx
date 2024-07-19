package middleware

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nskforward/httpx/cache"
	"github.com/nskforward/httpx/types"
)

func Cache(dir string, maxFileSize int64) types.Middleware {

	fi, err := os.Stat(dir)
	if err != nil {
		panic(fmt.Errorf("cache dir must be a valid path: %w", err))
	}
	if !fi.IsDir() {
		panic(fmt.Errorf("cache dir must be a directory"))
	}

	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if r.Method != http.MethodGet {
				return next(w, r)
			}

			if cacheSent(w, r) {
				return nil
			}

			var c *cache.Cache
			var f *os.File

			defer func() {
				if f != nil {
					f.Close()
				}
				if c != nil {
					cache.Set(r.URL.RequestURI(), c)
				}
			}()

			ww := fillCache(w, r, dir, maxFileSize, c, f)
			return next(ww, r)
		}
	}
}

func cacheSent(w http.ResponseWriter, r *http.Request) bool {
	c := cache.Get(r.URL.RequestURI())
	if c == nil {
		return false
	}
	if time.Since(c.To) > 0 {
		return false
	}
	err := c.SendFile(w, r)
	if err != nil {
		slog.Error("cannot send cache file", "error", err)
		return false
	}
	return true
}

func fillCache(w http.ResponseWriter, r *http.Request, dir string, maxFileSize int64, c *cache.Cache, f *os.File) *types.ResponseWrapper {
	ww := types.NewResponseWrapper(w)
	ww.BeforeBody = func() {
		if ww.Status() != 200 {
			return
		}

		cacheControlString := w.Header().Get(types.CacheControl)
		lastModifiedString := w.Header().Get(types.LastModified)
		etag := w.Header().Get(types.ETag)

		if cacheControlString == "" && lastModifiedString == "" && etag == "" {
			return
		}

		contentLengthString := w.Header().Get(types.ContentLength)
		if contentLengthString == "" {
			return
		}

		contentLength, err := strconv.ParseInt(contentLengthString, 10, 64)
		if err != nil {
			return
		}

		if contentLength > maxFileSize {
			return
		}

		control := cache.NewControl(cacheControlString)

		if control.NoStore || control.Private {
			return
		}

		if (w.Header().Get(types.SetCookie) != "" || r.Header.Get(types.Authorization) != "") && !control.Public {
			return
		}

		keyFolder := filepath.Join(dir, cache.Hash(r.URL.RequestURI()))
		os.Mkdir(keyFolder, os.ModePerm)
		keyFile := uuid.New().String()

		filename := filepath.Join(keyFolder, keyFile)
		f, err = os.Create(filename)
		if err != nil {
			slog.Error("cannot create cache file", "error", err, "file", filename)
			return
		}

		c = &cache.Cache{
			Key:         keyFolder,
			From:        cache.DetectDateMofified(w.Header()),
			To:          cache.DetectDateExpiration(w.Header(), control),
			Stale:       maxDuration(control.StaleWhileRevalidate, control.StaleIfError),
			ContentType: w.Header().Get(types.ContentType),
			ETag:        w.Header().Get(types.ETag),
		}

		encoding := detectEncoding(w.Header())
		if encoding == "gzip" {
			c.Filename.Gzip = keyFile
		} else {
			c.Filename.Plain = keyFile
		}

		ww.SetWriter(io.MultiWriter(f, ww.ResponseWriter))
	}

	return ww
}

func detectEncoding(header http.Header) string {
	if strings.Contains(header.Get(types.ContentEncoding), "gzip") {
		return "gzip"
	}
	return "plain"
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

/*
Usecases

// preventing caching anywhere
Cache-Control: no-store

// cache static content
Cache-Control: max-age=31536000, immutable

// allow cache, but up-to-date contents always
1. Cache-Control: no-cache
2. Cache-Control: max-age=0, must-revalidate
3. Last-Modified: Mon, 15 Jul 2024 13:37:33 GMT

w.Header().Set("Last-Modified", modtime.UTC().Format(TimeFormat))

func writeNotModified(w ResponseWriter) {
	// RFC 7232 section 4.1:
	// a sender SHOULD NOT generate representation metadata other than the
	// above listed fields unless said metadata exists for the purpose of
	// guiding cache updates (e.g., Last-Modified might be useful if the
	// response does not have an ETag field).
	h := w.Header()
	delete(h, "Content-Type")
	delete(h, "Content-Length")
	delete(h, "Content-Encoding")
	if h.Get("Etag") != "" {
		delete(h, "Last-Modified")
	}
	w.WriteHeader(StatusNotModified)
}

etag := header.Get("ETag")                             //  If-None-Match
lastModified := parseTime(header.Get("Last-Modified")) //  If-Modified-Since
// rw.ResponseWriter.Header().Get("Expires")
// rw.ResponseWriter.Header().Get("Age")
// rw.ResponseWriter.Header().Get("Date")
vary := header.Get("Vary")


HTTP/1.1 200 OK
Accept-Ranges: none
Age: 2919
Cache-Control: public, max-age=3600
Content-Encoding: br
Content-Length: 724
Content-Type: text/javascript
Date: Tue, 16 Jul 2024 07:09:54 GMT
ETag: W/"496e1e4ae72df3a4b647cb6bd577cf62"
Expires: Tue, 16 Jul 2024 07:14:09 GMT
Last-Modified: Mon, 15 Jul 2024 00:41:35 GMT
Server: Google Frontend
Strict-Transport-Security: max-age=63072000
Vary: Accept-Encoding
Via: 1.1 google
x-cache: hit
X-Content-Type-Options: nosniff
X-Frame-Options: DENY


GOLANG FILE SERVER RERSPONSE:
HTTP/1.1 200 OK
Content-Length: 195447
Accept-Ranges: bytes
Content-Type: text/html; charset=utf-8
Date: Thu, 18 Jul 2024 12:19:42 GMT
Last-Modified: Thu, 18 Jul 2024 12:15:27 GMT
*/
