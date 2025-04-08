package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ParseBody(r *http.Request, dst any) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return Throw(http.StatusUnsupportedMediaType, "Content-Type header must be present")
	}
	if strings.HasPrefix(contentType, "application/json") {
		return parseBodyJSON(r.Body, dst)
	}
	return Throw(http.StatusUnsupportedMediaType, fmt.Sprintf("unknown Content-Type '%s'", contentType))
}

func parseBodyJSON(r io.Reader, dst any) error {
	return json.NewDecoder(r).Decode(dst)
}
