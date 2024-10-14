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

func (s *Stream) sendField(field string, value []byte) error {
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
