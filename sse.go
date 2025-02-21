package httpx

import "time"

type sseSender func(
	send func(name, value string),
	flush func(),
) bool

// Stream returns FALSE if client gone, TRUE if server breaks stream.
func (ctx *Ctx) Stream(step sseSender) bool {
	ctx.SetHeader("Content-Type", "text/event-stream")
	ctx.SetHeader("Cache-Control", "no-store")

	gone := ctx.Request().Context().Done()
	writer := ctx.w
	flusher := ctx.w.Flusher()

	send := func(name, value string) {
		writer.Write([]byte(name))
		writer.Write([]byte(": "))
		writer.Write([]byte(value))
		writer.Write([]byte{'\n'})
	}

	flush := func() {
		writer.Write([]byte{'\n'})
		flusher.Flush()
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-gone:
			return false

		case <-ticker.C:
			writer.Write([]byte(":ping\n"))
			flusher.Flush()

		default:
			keepOpen := step(send, flush)
			if !keepOpen {
				return true
			}
		}
	}
}
