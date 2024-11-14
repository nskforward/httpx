package jwt

import "fmt"

type Token struct {
	Header    Header
	Claims    Claims
	signature string
}

func NewWithClaims(alg SignAlg, claims Claims) Token {
	return Token{
		Header: Header{
			A: alg,
			T: "JWT",
		},
		Claims: claims,
	}
}

func (t *Token) Sign(secret string) (string, error) {
	switch t.Header.A {
	case HS256:
		return signHS256(secret, t)

	default:
		return "", fmt.Errorf("jwt: unknown alg: %s", t.Header.A)
	}
}

func (t Token) Verify(secret string) error {
	originSignature := t.signature
	_, err := t.Sign(secret)
	if err != nil {
		return fmt.Errorf("jwt: bad token: %w", err)
	}
	if t.signature != originSignature {
		return fmt.Errorf("jwt: bad signature")
	}
	return nil
}
