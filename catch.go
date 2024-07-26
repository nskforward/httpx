package httpx

import (
	"net/http"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func Catch(handler types.Handler, logger types.LoggerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ww := types.NewResponseWrapper(w)

		err := handler(ww, r)

		if err != nil {
			handleError(ww, err)
		}

		if logger != nil {
			logger(ww, r, err)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	status := 400
	text := err.Error()

	apiError, ok := err.(response.APIError)
	if ok {
		status = apiError.Status
		if apiError.Text == "" || apiError.Status == 500 {
			text = http.StatusText(apiError.Status)
		} else {
			text = apiError.Text
		}
	}
	http.Error(w, text, status)
}
