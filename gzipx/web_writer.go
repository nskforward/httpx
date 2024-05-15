package gzipx

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"
)

const minContentLength = 2048

var allowedContentTypes = []string{
	"text/",
	"application/javascript",
	"application/json",
}

var bytesPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type WebWriter struct {
	status      int
	raw         http.ResponseWriter
	gz          *GZWriter
	buf         *bytes.Buffer
	headersSent bool
	ignore      bool
}

func NewWebWriter(raw http.ResponseWriter) *WebWriter {
	return &WebWriter{
		raw: raw,
		gz:  NewGZWriter(raw),
		buf: bytesPool.Get().(*bytes.Buffer),
	}
}

func (w *WebWriter) Close() {
	if !w.headersSent && w.buf.Len() > 0 {
		w.Header().Set("X-GZip-Ignore", "small content")
		w.raw.WriteHeader(200)
		io.Copy(w.raw, w.buf)
	} else {
		w.gz.Close()
	}
	w.buf.Reset()
	bytesPool.Put(w.buf)
}

func (w *WebWriter) Header() http.Header {
	return w.raw.Header()
}

func (w *WebWriter) check(status int) {
	w.status = status

	if status != http.StatusOK {
		w.ignore = true
		w.Header().Set("X-GZip-Ignore", "not ok")
		return
	}

	if w.Header().Get("Content-Encoding") != "" {
		w.ignore = true
		w.Header().Set("X-GZip-Ignore", "already compressed")
		return
	}

	if !IsSupportedContentType(w.Header().Get("Content-Type")) {
		w.ignore = true
		w.Header().Set("X-GZip-Ignore", "unsupported content")
		return
	}
}

func (w *WebWriter) WriteHeader(status int) {
	w.check(status)
	if w.ignore {
		w.raw.WriteHeader(status)
		w.headersSent = true
	}
}

func (w *WebWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(200)
	}

	// disabled compression
	if w.ignore {
		return w.raw.Write(b)
	}

	// compression enabled
	if w.headersSent {
		return w.gz.Write(b)
	}

	// try to fill internal buffer to decide if need compression
	if w.buf.Len()+len(b) < minContentLength {
		w.buf.Write(b)
		return len(b), nil
	}

	// buffer is filled up, start compression
	w.Header().Del("Content-Length")
	w.Header().Del("Accept-Ranges")
	w.Header().Set("Content-Encoding", "gzip")
	w.raw.WriteHeader(200)
	w.headersSent = true
	io.Copy(w.gz, w.buf) // flush internal buffer first and then write chank
	return w.gz.Write(b)
}

func IsSupportedContentType(contentType string) bool {
	if contentType == "" {
		return false
	}
	for _, ct := range allowedContentTypes {
		if strings.Contains(contentType, ct) {
			return true
		}
	}
	return false
}
