package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	ErrSSEventStreamClosed = fmt.Errorf("stream cloased")
)

type SSEventStream struct {
	w          http.ResponseWriter
	buf        bytes.Buffer
	controller *http.ResponseController
	encoder    *json.Encoder
	active     atomic.Bool
	intput     chan SSEvent
	ctx        context.Context
	cancel     context.CancelFunc
}

type SSEvent struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func NewSSEventStream(w http.ResponseWriter, r *http.Request) *SSEventStream {
	ctx, cancel := context.WithCancel(r.Context())

	s := &SSEventStream{
		w:          w,
		controller: http.NewResponseController(w),
		encoder:    json.NewEncoder(w),
		intput:     make(chan SSEvent, 4),
		ctx:        ctx,
		cancel:     cancel,
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")

	go s.loop()
	return s
}

func (s *SSEventStream) Close() {
	s.active.Store(false)
	s.cancel()
	close(s.intput)
}

func (s *SSEventStream) Next() bool {
	return s.active.Load()
}

func (s *SSEventStream) Send(event SSEvent) error {
	if !s.active.Load() {
		return ErrSSEventStreamClosed
	}
	s.intput <- event
	return nil
}

func (s *SSEventStream) loop() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	defer s.active.Store(false)

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-ticker.C:
			err := s.ping()
			if err != nil {
				return
			}

		case event := <-s.intput:
			err := s.write(event)
			if err != nil {
				return
			}
		}
	}
}

func (s *SSEventStream) ping() error {
	_, err := s.w.Write([]byte(":ping\n"))
	if err != nil {
		return err
	}
	return s.controller.Flush()
}

func (s *SSEventStream) write(event SSEvent) error {
	s.buf.WriteString("data: ")
	err := s.encoder.Encode(event)
	if err != nil {
		return err
	}
	s.buf.WriteString("\n\n")
	_, err = io.Copy(s.w, &s.buf)
	if err != nil {
		return err
	}
	return s.controller.Flush()
}
