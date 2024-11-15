package middleware

import (
	"fmt"
	"net/http"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

type BasicAuthOpts struct {
	Realm             string
	UsersAndPasswords map[string]string
	SkipFunc          func(r *http.Request) bool
	OnSuccess         func(user string, r *http.Request) *http.Request
}

func BasicAuth(opts BasicAuthOpts) types.Middleware {

	if len(opts.UsersAndPasswords) == 0 {
		panic(fmt.Errorf("basic auth middleware must contain at least one predefined user with password"))
	}

	realm := opts.Realm
	if realm == "" {
		realm = "Authentication required"
	}

	fail := func(w http.ResponseWriter) error {
		w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
		return response.APIError{Status: http.StatusUnauthorized, Text: "basic auth: incorrect user or password"}
	}

	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {

			if opts.SkipFunc != nil && opts.SkipFunc(r) {
				return next(w, r)
			}

			user, providedPass, ok := r.BasicAuth()
			if !ok {
				return fail(w)
			}

			rightPassword, ok := opts.UsersAndPasswords[user]
			if !ok {
				return fail(w)
			}

			if providedPass != rightPassword {
				return fail(w)
			}

			if opts.OnSuccess != nil {
				r = opts.OnSuccess(user, r)
			}

			return next(w, r)
		}
	}
}
