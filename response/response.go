package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/nskforward/httpx/types"
)

type H map[string]any

func NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func Redirect(w http.ResponseWriter, r *http.Request, status int, url string) error {
	http.Redirect(w, r, url, status)
	return nil
}

func JSON(w http.ResponseWriter, statusCode int, obj any) error {
	w.Header().Set(types.ContentType, types.ApplicationJSON)
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(obj)
}

func Text(w http.ResponseWriter, statusCode int, msg string) error {
	w.Header().Set(types.ContentType, types.TextPlain)
	w.Header().Set(types.ContentLength, strconv.Itoa(len(msg)))
	w.WriteHeader(statusCode)
	_, err := io.WriteString(w, msg)
	return err
}

func RawData(w http.ResponseWriter, statusCode int, contentType string, body []byte) error {
	w.Header().Set(types.ContentType, contentType)
	w.Header().Set(types.ContentLength, strconv.Itoa(len(body)))
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	return err
}

func Range(w http.ResponseWriter, statusCode int, contentType string, r io.ReadSeeker) error {
	// ranged response
	return fmt.Errorf("not implemented")
}
