package httpx

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type File struct {
	mx        sync.RWMutex
	filepath  string
	body      []byte
	length    string
	headers   map[string]string
	stat      fs.FileInfo
	lastCheck time.Time
}

func NewFile(filepath string, headers map[string]string) (*File, error) {
	initialStat, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	body, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return &File{
		filepath:  filepath,
		body:      body,
		length:    strconv.Itoa(len(body)),
		headers:   headers,
		stat:      initialStat,
		lastCheck: time.Now(),
	}, nil
}

func (f *File) Serve(w http.ResponseWriter, r *http.Request) {
	err := f.check()
	if err != nil {
		fmt.Println("error: cannot serve file:", f.filepath, "-", err)
		http.Error(w, "cannot serve file", 404)
		return
	}
	f.flush(w)
}

func (f *File) flush(w http.ResponseWriter) {
	for k, v := range f.headers {
		w.Header().Set(k, v)
	}
	f.mx.RLock()
	defer f.mx.RUnlock()

	w.Header().Set("Content-Length", f.length)
	w.WriteHeader(200)
	w.Write(f.body)
}

func (f *File) check() error {
	f.mx.RLock()
	isFresh := time.Since(f.lastCheck) < time.Second
	f.mx.RUnlock()

	if isFresh {
		return nil
	}

	stat, err := os.Stat(f.filepath)
	if err != nil {
		return err
	}

	f.mx.Lock()
	f.lastCheck = time.Now()
	changed := stat.Size() != f.stat.Size() || stat.ModTime() != f.stat.ModTime()
	f.mx.Unlock()

	if !changed {
		return nil
	}

	err = f.update(stat)
	return err
}

func (f *File) update(stat fs.FileInfo) error {
	f.mx.Lock()
	defer f.mx.Unlock()

	f.stat = stat
	body, err := os.ReadFile(f.filepath)
	if err != nil {
		return err
	}
	f.body = body
	f.length = strconv.Itoa(len(body))
	return nil
}
