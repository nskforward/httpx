package jwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

type SignAlg string

const (
	HS256 SignAlg = "HS256"
)

func signHS256(secret string, t *Token) (string, error) {
	h, err := json.Marshal(t.Header)
	if err != nil {
		return "", err
	}

	c, err := json.Marshal(t.Claims)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.WriteString(base64.StdEncoding.EncodeToString(h))
	buf.WriteByte('.')
	buf.WriteString(base64.RawURLEncoding.EncodeToString(c))

	hasher := hmac.New(sha256.New, []byte(secret))
	hasher.Write(buf.Bytes())

	signature := hasher.Sum(nil)
	t.signature = base64.RawURLEncoding.EncodeToString(signature)

	buf.WriteByte('.')
	buf.WriteString(t.signature)

	return buf.String(), nil
}
