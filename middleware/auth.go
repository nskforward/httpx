package middleware

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/nskforward/httpx/jwt"
	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

var ContextAuthString types.ContextParam = "middleware.auth.string"

func JWTAuth(secret string) types.Middleware {

	encoder := jwt.NewEncoder(secret)

	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			token := r.Header.Get(types.Authorization)
			if token == "" {
				return response.APIError{Status: http.StatusUnauthorized, Text: "require Authorization header"}
			}
			data, err := encoder.Decode([]byte(token))
			if err != nil {
				return response.APIError{Status: http.StatusUnauthorized, Text: "bad Authorization header"}
			}
			r = types.SetParam(r, ContextAuthString, string(data))
			return next(w, r)
		}
	}
}

func BasicAuth(realm string, creds map[string]string) types.Middleware {
	basicAuthFailed := func(w http.ResponseWriter, realm string) error {
		w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
		return response.NewAPIError(http.StatusUnauthorized)
	}
	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			user, pass, ok := r.BasicAuth()
			if !ok {
				return basicAuthFailed(w, realm)
			}
			credPass, ok := creds[user]
			if !ok || subtle.ConstantTimeCompare([]byte(pass), []byte(credPass)) != 1 {
				return basicAuthFailed(w, realm)
			}
			r = types.SetParam(r, ContextAuthString, user)
			return next(w, r)
		}
	}
}

func GetAuthString(r *http.Request) string {
	val := types.GetParam(r, ContextAuthString)
	if val == nil {
		return ""
	}
	return val.(string)
}
