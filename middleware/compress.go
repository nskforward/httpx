package middleware

import (
	"compress/gzip"
	"net/http"
	"strconv"
	"strings"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func Compress(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if !strings.Contains(r.Header.Get(types.AcceptEncoding), "gzip") {
			return next(w, r)
		}

		ww := response.NewWrapper(w)
		var gzw *gzip.Writer
		defer func() {
			if gzw != nil {
				gzw.Close()
			}
		}()

		ww.BeforeBody = func() {
			if !isAllowedContent(ww.Header().Get(types.ContentType)) {
				return
			}
			if !isAllowedLength(ww.Header().Get(types.ContentLength)) {
				return
			}
			gz, err := gzip.NewWriterLevel(ww.ResponseWriter, 6)
			if err != nil {
				return
			}
			gzw = gz
			ww.Header().Del(types.ContentLength)
			ww.Header().Del(types.AcceptRanges)
			ww.Header().Set(types.ContentEncoding, "gzip")
			ww.BodyWriter(gz)
		}

		return next(ww, r)
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
	return size > 2048
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
