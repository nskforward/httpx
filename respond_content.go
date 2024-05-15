package httpx

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

const (
	ContentTypeJSON = "application/json; charset=utf-8"
	ContentTypeHTML = "text/html; charset=utf-8"
	ContentTypeTEXT = "text/plain; charset=utf-8"
)

func NoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func JSON(w http.ResponseWriter, statusCode int, obj any) error {
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(obj)
}

func Text(w http.ResponseWriter, statusCode int, msg string) error {
	w.Header().Set("Content-Type", ContentTypeTEXT)
	w.Header().Set("Content-Length", strconv.Itoa(len(msg)))
	w.WriteHeader(statusCode)
	_, err := io.WriteString(w, msg)
	return err
}

func Data(w http.ResponseWriter, statusCode int, contentType string, body []byte) error {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	return err
}
