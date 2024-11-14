package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func DecodeBody(r *http.Request, dest any) error {
	contentType := r.Header.Get(types.ContentType)
	if contentType == "" {
		return response.APIError{Status: http.StatusUnsupportedMediaType, Text: "Unknown Media Type"}
	}
	if strings.HasPrefix(contentType, "application/json") {
		return decodeJSON(r.Body, dest)
	}
	return response.APIError{Status: http.StatusUnsupportedMediaType, Text: fmt.Sprintf("Unsupported Media Type: %s", contentType)}
}

func decodeJSON(r io.Reader, dest any) error {
	return json.NewDecoder(r).Decode(dest)
}
