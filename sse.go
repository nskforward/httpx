package httpx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Stream struct {
	ctx     context.Context
	cancel  context.CancelFunc
	flusher http.Flusher
	output  io.Writer
	queue   chan *StreamEvent
}

type StreamEvent struct {
	conn *Stream
	name string
	id   string
	buf  *bytes.Buffer
}

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 512))
	},
}

func NewStream(w http.ResponseWriter, r *http.Request) (*Stream, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("http.ResponseWriter instance must implemend http.Flusher")
	}
	ctx, cancel := context.WithCancel(r.Context())

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")

	output := &Stream{output: w, flusher: flusher, ctx: ctx, cancel: cancel, queue: make(chan *StreamEvent, 32)}
	go output.handleQueue()

	return output, nil
}

func (s *Stream) Close() {
	s.cancel()
	close(s.queue)
}

func (s *Stream) Alive() bool {
	return s.ctx.Err() == nil
}

func (s *Stream) handleQueue() {
	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.output.Write([]byte(":ping\n"))
		case event, ok := <-s.queue:
			if !ok {
				return
			}
			s.send(event)
		}
	}
}

func (s *Stream) send(event *StreamEvent) {
	defer bufferPool.Put(event.buf)

	if event.name != "" {
		s.output.Write([]byte("event: "))
		s.output.Write([]byte(event.name))
		s.output.Write([]byte("\n"))
	}

	if event.id != "" {
		s.output.Write([]byte("id: "))
		s.output.Write([]byte(event.id))
		s.output.Write([]byte("\n"))
	}

	_, err := s.output.Write([]byte("data: "))
	if err != nil {
		return
	}
	io.Copy(s.output, event.buf)
	s.output.Write([]byte("\n\n"))
}

func (s *Stream) Event() StreamEvent {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return StreamEvent{
		conn: s,
		buf:  buf,
	}
}

func (event *StreamEvent) SetName(name string) {
	event.name = name
}

func (event *StreamEvent) Context() context.Context {
	return event.conn.ctx
}

func (event *StreamEvent) SetID(id string) {
	event.id = id
}

func (event *StreamEvent) Write(p []byte) (int, error) {
	return event.buf.Write(p)
}

func (event *StreamEvent) WriteString(s string) (int, error) {
	return event.buf.WriteString(s)
}

func (event *StreamEvent) Send() {
	select {
	case event.conn.queue <- event:
	default:
	}
}
