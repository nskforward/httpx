package httpx

import (
	"net/http"
	"slices"
	"strings"

	"github.com/nskforward/httpx/types"
)

func ParseAuthToken(r *http.Request) (token, scheme string, ok bool) {
	authorization := r.Header.Get(types.Authorization)
	if authorization == "" {
		return "", "", false
	}
	items := strings.Split(authorization, " ")
	items = slices.DeleteFunc(items, func(item string) bool {
		return item == " "
	})
	if len(items) == 2 {
		return items[1], items[0], true
	}
	if len(items) == 1 {
		return items[0], "", true
	}
	return "", "", false
}
