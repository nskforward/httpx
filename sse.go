package httpx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
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

func SSE(w http.ResponseWriter, r *http.Request, callback ...func(s *Stream) bool) error {
	s, err := newStream(w, r)
	if err != nil {
		return err
	}
	defer s.Close()

	for _, f := range callback {
		select {
		case <-r.Context().Done():
			return nil
		default:
			if !f(s) {
				return nil
			}
		}
	}

	return nil
}

func newStream(w http.ResponseWriter, r *http.Request) (*Stream, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("http.ResponseWriter instance must implemend http.Flusher")
	}

	fmt.Println("stream supports flush")

	ctx, cancel := context.WithCancel(r.Context())

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")

	output := &Stream{output: w, flusher: flusher, ctx: ctx, cancel: cancel, queue: make(chan *StreamEvent, 32)}
	go output.handleQueue()

	return output, nil
}

func (s *Stream) Close() {
	close(s.queue)
}

func (s *Stream) Context() context.Context {
	return s.ctx
}

func (s *Stream) handleQueue() {
	ticker := time.NewTicker(15 * time.Second)
	for {
		select {

		case <-s.ctx.Done():
			slog.Info("SSE context done")
			return

		case <-ticker.C:
			s.output.Write([]byte(":ping\n"))

		case event, ok := <-s.queue:
			if !ok {
				slog.Info("SSE dequeued empty")
				s.cancel()
				return
			}
			err := s.send(event)
			if err != nil {
				slog.Info("SSE send failed", "error", err)
				s.cancel()
				return
			}
		}
	}
}

func (s *Stream) send(event *StreamEvent) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)
	defer bufferPool.Put(event.buf)

	if event.name != "" {
		s.buffFillString(buf, "event", event.name)
	}

	if event.id != "" {
		s.buffFillString(buf, "id", event.id)
	}

	s.buffFillBytes(buf, "data", event.buf.Bytes())

	buf.WriteByte('\n')

	_, err := io.Copy(s.output, buf)
	if err != nil {
		return err
	}
	s.flusher.Flush()

	return nil
}

func (s *Stream) buffFillBytes(buf *bytes.Buffer, field string, value []byte) {
	buf.WriteString(field)
	buf.WriteString(": ")
	buf.Write(value)
	buf.WriteByte('\n')
}

func (s *Stream) buffFillString(buf *bytes.Buffer, field, value string) {
	buf.WriteString(field)
	buf.WriteString(": ")
	buf.WriteString(value)
	buf.WriteByte('\n')
}

func (s *Stream) WriteString(msg string) {
	event := s.Event()
	event.buf.WriteString(msg)
	event.Send()
}

func (s *Stream) Event() *StreamEvent {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return &StreamEvent{
		conn: s,
		buf:  buf,
	}
}

func (event *StreamEvent) SetName(name string) {
	event.name = name
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
		slog.Info("SSE enqueued")
	default:
		slog.Info("SSE not enqueued")
	}
}
