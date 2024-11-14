package auth

import (
	"encoding/base64"
	"fmt"
	"strings"
)

var _ Authenticator = (*BasicAuthenticator)(nil)

type BasicAuthenticator struct {
	users map[string]string
}

func NewBasicAuthenticator(users map[string]string) *BasicAuthenticator {
	return &BasicAuthenticator{
		users: users,
	}
}

func (a *BasicAuthenticator) Authenticate(token string) (any, error) {
	c, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("basic auth: %w", err)
	}
	cs := string(c)
	username, password, ok := strings.Cut(cs, ":")
	if !ok {
		return nil, fmt.Errorf("basic auth: bad token format")
	}
	storedPass, ok := a.users[username]
	if !ok {
		return nil, fmt.Errorf("basic auth: unknown user: %s", username)
	}
	if !strings.EqualFold(storedPass, password) {
		return nil, fmt.Errorf("basic auth: incorrect password for user: %s", username)
	}
	return username, nil
}
