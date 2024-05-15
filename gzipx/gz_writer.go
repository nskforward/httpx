package gzipx

import (
	"compress/gzip"
	"io"
	"sync"
)

var gzPool = sync.Pool{
	New: func() interface{} {
		gz, err := gzip.NewWriterLevel(io.Discard, 6)
		if err != nil {
			panic(err)
		}
		return gz
	},
}

type GZWriter struct {
	zw *gzip.Writer
}

func NewGZWriter(w io.Writer) *GZWriter {
	zw := gzPool.Get().(*gzip.Writer)
	zw.Reset(w)
	return &GZWriter{zw}
}

func (w *GZWriter) Close() {
	w.zw.Flush()
	w.zw.Close()
	gzPool.Put(w.zw)
}

func (w *GZWriter) Write(b []byte) (int, error) {
	return w.zw.Write(b)
}
