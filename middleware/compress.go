package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

var gzPool = sync.Pool{
	New: func() any {
		return gzip.NewWriter(io.Discard)
	},
}

func Compress(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if !strings.Contains(r.Header.Get(types.AcceptEncoding), "gzip") {
			return next(w, r)
		}

		ww := response.NewWrapper(w)
		var gz *gzip.Writer

		ww.BeforeBody = func() {
			if !isAllowedContent(ww.Header().Get(types.ContentType)) {
				return
			}
			if !isAllowedLength(ww.Header().Get(types.ContentLength)) {
				return
			}

			gz = gzPool.Get().(*gzip.Writer)
			gz.Reset(w)

			ww.SetWriter(gz)
			ww.Header().Del(types.ContentLength)
			ww.Header().Del(types.AcceptRanges)
			w.Header().Set(types.ContentEncoding, "gzip")
		}

		err := next(ww, r)

		if gz != nil {
			gz.Close()
			gzPool.Put(gz)
		}

		return err
	}
}

func isAllowedLength(contentLength string) bool {
	if contentLength == "" {
		return false
	}
	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return false
	}
	return size > 2000
}

func isAllowedContent(contentType string) bool {
	if contentType == "" {
		return false
	}
	for _, part := range []string{"text/", "application/javascript", "application/json"} {
		if strings.Contains(contentType, part) {
			return true
		}
	}
	return false
}
