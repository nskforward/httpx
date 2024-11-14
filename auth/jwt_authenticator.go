package auth

import (
	"fmt"
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
	return nil, fmt.Errorf("not implemented")
}
