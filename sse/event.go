package sse

type Event struct {
	Name string
	ID   string
	Data []byte
}
