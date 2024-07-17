package cache

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

type KeyFunc func(r *http.Request) string

var (
	keyFuncStoreURI sync.Map // store key func by uri
	updatedStoreKey sync.Map // store last update by key
)

func CreateETag(r *http.Request, keyFunc KeyFunc) (etag string, updated time.Time) {
	uri := r.URL.RequestURI()
	if _, ok := keyFuncStoreURI.Load(uri); !ok {
		keyFuncStoreURI.Store(r.URL.RequestURI(), keyFunc)
	}
	key := keyFunc(r)
	updated = time.Now()
	updatedStoreKey.Store(key, updated)
	etag = generateETag(key, updated)
	return
}

func SearchETag(r *http.Request) (etag string, updated time.Time) {
	uri := r.URL.RequestURI()
	keyFunc, ok := keyFuncStoreURI.Load(uri)
	if !ok {
		return
	}
	key := keyFunc.(KeyFunc)(r)
	up, ok := updatedStoreKey.Load(key)
	if !ok {
		return
	}
	updated = up.(time.Time)
	etag = generateETag(key, updated)
	return
}

func generateETag(key string, updated time.Time) string {
	return Hash(key, strconv.FormatInt(updated.UnixNano(), 36))
}
