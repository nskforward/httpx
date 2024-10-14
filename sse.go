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

type stream struct {
	ctx     context.Context
	cancel  context.CancelFunc
	flusher http.Flusher
	output  io.Writer
	queue   chan *StreamEvent
}

type StreamEvent struct {
	conn *stream
	name string
	id   string
	buf  *bytes.Buffer
}

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 512))
	},
}

func Stream(w http.ResponseWriter, r *http.Request, callback ...func(s *stream) bool) error {
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

func newStream(w http.ResponseWriter, r *http.Request) (*stream, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("http.ResponseWriter instance must implemend http.Flusher")
	}
	ctx, cancel := context.WithCancel(r.Context())

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")

	w.WriteHeader(http.StatusNoContent)

	output := &stream{output: w, flusher: flusher, ctx: ctx, cancel: cancel, queue: make(chan *StreamEvent, 32)}
	go output.handleQueue()

	return output, nil
}

func (s *stream) Close() {
	close(s.queue)
}

func (s *stream) Context() context.Context {
	return s.ctx
}

func (s *stream) handleQueue() {
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

func (s *stream) send(event *StreamEvent) error {
	defer bufferPool.Put(event.buf)

	err := s.sendField("event", []byte(event.name))
	if err != nil {
		return err
	}

	err = s.sendField("id", []byte(event.id))
	if err != nil {
		return err
	}

	err = s.sendField("data", event.buf.Bytes())
	if err != nil {
		return err
	}

	_, err = s.output.Write([]byte("\n"))
	return err
}

func (s *stream) sendField(field string, value []byte) error {
	if len(value) == 0 {
		slog.Info("SSE send with empty value", "field", field)
		return nil
	}

	_, err := s.output.Write([]byte(field))
	if err != nil {
		return err
	}

	_, err = s.output.Write([]byte(": "))
	if err != nil {
		return err
	}
	_, err = s.output.Write([]byte(value))
	if err != nil {
		return err
	}
	_, err = s.output.Write([]byte("\n"))
	if err != nil {
		return err
	}

	slog.Info("SSE send", "field", field, "value", string(value))

	return nil
}

func (s *stream) WriteString(msg string) {
	event := s.Event()
	event.buf.WriteString(msg)
	event.Send()
}

func (s *stream) Event() *StreamEvent {
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
