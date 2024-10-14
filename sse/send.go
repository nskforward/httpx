package sse

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 512))
	},
}

func send(flusher http.Flusher, w http.ResponseWriter, e Event) {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	if e.Name != "" {
		writeField(buf, "event", []byte(e.Name))
	}

	if e.ID != "" {
		writeField(buf, "id", []byte(e.ID))
	}

	if len(e.Data) > 0 {
		writeField(buf, "data", e.Data)
	}

	if buf.Len() > 0 {
		buf.WriteByte('\n')
	}

	io.Copy(w, buf)
	flusher.Flush()
}

func writeField(buf *bytes.Buffer, field string, value []byte) {
	buf.WriteString(field)
	buf.WriteString(": ")
	buf.Write(value)
	buf.WriteByte('\n')
}
