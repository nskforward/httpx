package logging

import "sync"

var pool = sync.Pool{
	New: func() interface{} {
		return new(Writer)
	},
}

func Get() *Writer {
	return pool.Get().(*Writer)
}

func Put(w *Writer) {
	pool.Put(w)
}
