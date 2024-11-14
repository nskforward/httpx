package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func Parse(s string) (*Token, error) {
	items := strings.Split(s, ".")
	if len(items) != 3 {
		return nil, fmt.Errorf("jwt: bad token format")
	}
	t := Token{
		signature: items[2],
	}

	data, err := base64.RawURLEncoding.DecodeString(items[0])
	if err != nil {
		return nil, fmt.Errorf("jwt: cannot decode token header: %w", err)
	}

	err = json.Unmarshal(data, &t.Header)
	if err != nil {
		return nil, fmt.Errorf("jwt: cannot decode token header: %w", err)
	}

	data, err = base64.RawURLEncoding.DecodeString(items[1])
	if err != nil {
		return nil, fmt.Errorf("jwt: cannot decode token payload: %w", err)
	}

	err = json.Unmarshal(data, &t.Claims)
	if err != nil {
		return nil, fmt.Errorf("jwt: cannot decode token payload: %w", err)
	}

	return &t, nil
}
