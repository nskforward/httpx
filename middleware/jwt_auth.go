package middleware

import (
	"fmt"
	"net/http"

	"github.com/nskforward/httpx"
	"github.com/nskforward/httpx/jwt"
	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

type JWTAuthOpts struct {
	Secret    string
	SkipFunc  func(r *http.Request) bool
	OnSuccess func(r *http.Request, claims jwt.Claims) *http.Request
}

func JWTAuth(opts JWTAuthOpts) types.Middleware {

	if len(opts.Secret) < 8 {
		panic(fmt.Errorf("jwt secret must be at least 8 chars length"))
	}

	fail := func() error {
		return response.APIError{Status: http.StatusUnauthorized, Text: "bad authorization token"}
	}

	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if opts.SkipFunc != nil && opts.SkipFunc(r) {
				return next(w, r)
			}
			token, _, ok := httpx.ParseAuthToken(r)
			if !ok {
				return fail()
			}
			t, err := jwt.Parse(token)
			if err != nil {
				return fail()
			}
			err = t.Verify(opts.Secret)
			if err != nil {
				return fail()
			}
			if opts.OnSuccess != nil {
				r = opts.OnSuccess(r, t.Claims)
			}
			return next(w, r)
		}
	}
}
