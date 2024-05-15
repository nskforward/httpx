package httpx

import (
	"net/http"
	"strings"

	"github.com/nskforward/httpx/gzipx"
)

func GZip(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("X-GZip-Ignore", "unsupported client")
			return next(w, r)
		}

		ww := gzipx.NewWebWriter(w)
		defer ww.Close()

		err := next(ww, r)

		if err != nil {
			return InternalServerError(err)
		}

		return err
	}
}
