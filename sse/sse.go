package sse

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func Stream(w http.ResponseWriter, r *http.Request, content ...func(ctx context.Context, w Writer) bool) error {

	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("stream unsupported")
	}

	queue := make(chan Event, 64)
	defer close(queue)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-r.Context().Done():
				return
			case <-ticker.C:
				w.Write([]byte(":ping\n"))
				flusher.Flush()
			case e := <-queue:
				if len(e.Data) == 0 {
					wg.Done()
					return
				}
				send(flusher, w, e)
			}
		}
	}()

	for _, f := range content {
		if !f(r.Context(), Writer{queue}) {
			break
		}
	}

	queue <- Event{}
	wg.Wait()

	return nil
}
