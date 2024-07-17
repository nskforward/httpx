package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func Decode(token string, secret []byte) ([]byte, error) {
	items := strings.Split(token, ".")
	if len(items) != 2 {
		return nil, fmt.Errorf("incorrect format")
	}
	payload, err := hex.DecodeString(items[0])
	if err != nil {
		return nil, fmt.Errorf("incorrect format")
	}
	signature := hex.EncodeToString(hmac.New(sha256.New, secret).Sum(payload))
	if !strings.EqualFold(signature, items[1]) {
		return nil, fmt.Errorf("invalid signature")
	}
	return payload, nil
}
