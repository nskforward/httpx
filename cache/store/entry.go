package store

import (
	"net/http"
	"time"
)

type Entry struct {
	FullPath   string
	Header     map[string]string
	Expiration time.Time
}

func (entry *Entry) Expired() bool {
	return time.Since(entry.Expiration) > 0
}

func (entry *Entry) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range entry.Header {
		if v != "" {
			w.Header().Set(k, v)
		}
	}
	http.ServeFile(w, r, entry.FullPath)
}
