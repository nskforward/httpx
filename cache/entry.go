package cache

import (
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/nskforward/httpx/types"
)

type Entry struct {
	sync.RWMutex
	bucket  *Bucket
	deleted atomic.Bool
	id      string
	key     string
	file    string
	status  uint32
	used    uint64
	from    time.Time
	to      time.Time
}

const (
	idle uint32 = iota
	filling
)

func (e *Entry) Valid() bool {
	return time.Until(e.to) > 0
}

func (e *Entry) ID() string {
	return e.id
}

func (e *Entry) File() string {
	return e.file
}

func (e *Entry) From() time.Time {
	return e.from
}

func (e *Entry) ChangeID() {
	e.id = uuid.New().String()
}

func (e *Entry) SetIdle() {
	atomic.StoreUint32(&e.status, idle)
}

func (e *Entry) IsIdle() bool {
	return atomic.LoadUint32(&e.status) == idle
}

func (e *Entry) SetFilling() bool {
	return atomic.CompareAndSwapUint32(&e.status, idle, filling)
}

func (e *Entry) Key() string {
	return e.key
}

func (e *Entry) Used() uint64 {
	return atomic.LoadUint64(&e.used)
}
func (e *Entry) Hit() {
	atomic.AddUint64(&e.used, 1)
}

func (e *Entry) IsDeleted() bool {
	return e.deleted.Load()
}

func (e *Entry) Delete() {
	e.Lock()
	defer e.Unlock()
	e.deleted.Store(true)
	e.bucket.keyStore.Delete(e.key)
	os.RemoveAll(e.file)
}

func (e *Entry) LastModified() string {
	return e.from.UTC().Format(http.TimeFormat)
}

func (e *Entry) SendNoCache(w http.ResponseWriter) {
	w.Header().Set(types.CacheControl, "no-cache")
	w.Header().Set(types.ETag, e.ID())
	w.Header().Set(types.LastModified, e.LastModified())
}
