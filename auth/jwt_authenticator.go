package auth

import (
	"github.com/nskforward/httpx/jwt"
)

var _ Authenticator = (*BasicAuthenticator)(nil)

type JWTAuthenticator struct {
	secret string
}

func NewJWTAuthenticator(secret string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
	}
}

func (a *JWTAuthenticator) Authenticate(token string) (any, error) {
	t, err := jwt.Parse(token)
	if err != nil {
		return nil, err
	}
	err = t.Verify(a.secret)
	if err != nil {
		return nil, err
	}
	return t.Claims, nil
}
