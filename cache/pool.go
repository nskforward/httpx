package cache

import (
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return new(Writer)
	},
}

func GetWriter() *Writer {
	return pool.Get().(*Writer)
}

func PutWriter(w *Writer) {
	pool.Put(w)
}
