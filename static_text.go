package httpx

import (
	"net/http"

	"github.com/nskforward/httpx/types"
)

func Text(text string) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte(text))
		return nil
	}
}
