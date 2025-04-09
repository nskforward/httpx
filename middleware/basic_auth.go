package middleware

import (
	"net/http"

	"github.com/nskforward/httpx"
	"golang.org/x/crypto/bcrypt"
)

// creds input param accept a map in format username=password where password string must be encoded by bcrypt
func BasicAuth(creds map[string]string) httpx.Handler {
	return func(req *http.Request, resp *httpx.Response) error {
		user, pass, ok := req.BasicAuth()
		if ok {
			hashedPassword, ok := creds[user]
			if ok {
				if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pass)) == nil {
					return resp.Next(req)
				}
			}
		}
		return resp.Unauthorized("require valid credentials")
	}
}
