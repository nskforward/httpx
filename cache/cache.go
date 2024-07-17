package cache

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nskforward/httpx/types"
)

var cacheStoreURI sync.Map // map[ContentEncoding]*Cache

type Cache struct {
	Key         string
	ContentType string
	From        time.Time
	To          time.Time
	Stale       time.Duration
	ETag        string
	Filename    struct {
		Plain string
		Gzip  string
	}
}

func (c *Cache) SendFile(w http.ResponseWriter, r *http.Request) error {
	var filename string

	allowGzip := c.Filename.Gzip != "" && strings.Contains(r.Header.Get(types.AcceptEncoding), "gzip")

	if allowGzip {
		filename = filepath.Join(c.Key, c.Filename.Gzip)
	} else {
		filename = filepath.Join(c.Key, c.Filename.Plain)
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}

	w.Header().Set(types.AcceptRanges, "none")
	w.Header().Set(types.Vary, "Accept-Encoding")
	w.Header().Set(types.XCache, "hit")
	w.Header().Set(types.ContentType, c.ContentType)
	w.Header().Set(types.Age, fmt.Sprintf("%.0f", time.Since(c.From).Seconds()))
	w.Header().Set(types.CacheControl, fmt.Sprintf("public, max-age=%.0f", c.To.Sub(c.From).Seconds()))
	w.Header().Set(types.ContentLength, fmt.Sprintf("%d", fi.Size()))
	w.Header().Set(types.Expires, c.To.UTC().Format(http.TimeFormat))
	w.Header().Set(types.LastModified, c.From.UTC().Format(http.TimeFormat))
	if allowGzip {
		w.Header().Set(types.ContentEncoding, "gzip")
	}
	if c.ETag != "" {
		w.Header().Set(types.ETag, c.ETag)
	}

	io.Copy(w, f)
	return nil
}

func Get(uri string) *Cache {
	c, ok := cacheStoreURI.Load(uri)
	if !ok {
		return nil
	}
	return c.(*Cache)
}

func Set(uri string, c *Cache) {
	old, ok := cacheStoreURI.Load(uri)
	if ok {
		Del(uri, old.(*Cache))
	}
	cacheStoreURI.Store(uri, c)
}

func Del(uri string, c *Cache) {
	cacheStoreURI.Delete(uri)
	if c != nil && c.Key != "" {
		os.RemoveAll(c.Key)
	}
}
