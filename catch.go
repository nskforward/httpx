package httpx

import (
	"net/http"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

func (ro *Router) Catch(next types.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = traceIDSetter(ro.logging(next))(w, r)
	}
}

func (ro *Router) logging(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ww := types.NewResponseWrapper(w)
		err := next(ww, r)
		if err != nil {
			ro.handleError(ww, r, err)
		}
		ro.loggerFunc(ww, r)
		return nil
	}
}

func (ro *Router) handleError(w http.ResponseWriter, r *http.Request, err error) {
	status := 400
	text := err.Error()

	apiError, ok := err.(response.APIError)
	if ok {
		status = apiError.Status
		if apiError.Text == "" {
			text = http.StatusText(apiError.Status)
		} else {
			text = apiError.Text
		}
	}

	ro.errorFunc(w, r, status, text)
}
