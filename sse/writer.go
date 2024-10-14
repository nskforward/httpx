package sse

type Writer struct {
	queue chan<- Event
}

func (w Writer) WriteString(msg string) {
	w.Write(Event{
		Data: []byte(msg),
	})
}

func (w Writer) Write(event Event) {
	w.queue <- event
}
