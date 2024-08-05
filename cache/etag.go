package cache

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type KeyFunc func(*http.Request) string

var keyFuncs sync.Map
var etags sync.Map

func LoadEtag(r *http.Request) string {
	f, ok := keyFuncs.Load(r.URL.RequestURI())
	if !ok {
		return ""
	}
	key := f.(KeyFunc)(r)
	if key == "" {
		return ""
	}
	etag, ok := etags.Load(key)
	if !ok {
		return ""
	}
	return etag.(string)
}

func SaveEtag(r *http.Request, keyFunc KeyFunc) string {
	key := keyFunc(r)
	if key == "" {
		return ""
	}
	etag := fmt.Sprintf("W/\"%s\"", strconv.FormatInt(time.Now().UnixNano(), 36))
	etags.Store(key, etag)
	keyFuncs.Store(r.URL.Path, keyFunc)
	return etag
}
