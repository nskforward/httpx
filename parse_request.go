package httpx

import (
	"encoding/json"
	"net/http"
)

func ParseRequest(r *http.Request, obj any) error {
	return json.NewDecoder(r.Body).Decode(obj)
}
