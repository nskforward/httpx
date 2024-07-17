package jwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Encode(secret, payload []byte) string {
	var buf bytes.Buffer
	buf.WriteString(hex.EncodeToString(payload))
	buf.WriteByte('.')
	buf.WriteString(hex.EncodeToString(hmac.New(sha256.New, []byte(secret)).Sum(payload)))
	return buf.String()
}
