package cache

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
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

func (e *Entry) From() time.Time {
	return e.from
}

func (e *Entry) File() string {
	return e.file
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

func (e *Entry) SendCache(w http.ResponseWriter) error {
	e.RLock()
	defer e.RUnlock()
	f, err := os.Open(e.file)
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
	w.Header().Set("Age", strconv.FormatFloat(time.Since(e.from).Seconds(), 'f', 0, 64))
	io.Copy(w, b)
	return nil
}

func (e *Entry) Delete() {
	e.Lock()
	defer e.Unlock()
	e.deleted.Store(true)
	e.bucket.keyStore.Delete(e.key)
	os.RemoveAll(e.file)
}
