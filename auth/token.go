package auth

import (
	"net/http"
	"slices"
	"strings"

	"github.com/nskforward/httpx/types"
)

func ParseToken(r *http.Request) (ok bool, scheme, token string) {
	authorization := r.Header.Get(types.Authorization)
	if authorization == "" {
		return false, "", ""
	}
	items := strings.Split(authorization, " ")
	items = slices.DeleteFunc(items, func(item string) bool {
		return item == " "
	})
	if len(items) == 2 {
		return true, items[0], items[1]
	}
	if len(items) == 1 {
		return true, "", items[0]
	}
	return false, "", ""
}
