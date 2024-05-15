package httpx

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type contextKey string

var userIDKey contextKey = "UserID"

func HandlerAuthJWT(secret string) func(next Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			token := r.Header.Get("Authorization")
			if token == "" {
				return Unauthorized(fmt.Errorf("authorization header is absent"))
			}

			data, err := JWTDecode(secret, token)
			if err != nil {
				return Unauthorized(fmt.Errorf("cannot decode the authorization token"))
			}

			ctx := context.WithValue(r.Context(), userIDKey, string(data))
			*r = *r.WithContext(ctx)
			err = next(w, r)
			return err
		}
	}
}

func UserID(r *http.Request) string {
	v := r.Context().Value(userIDKey)
	if v == nil {
		return ""
	}
	return v.(string)
}

var jwtBufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func JWTEncode(secret, payload string) string {
	buf := jwtBufferPool.Get().(*bytes.Buffer)
	defer jwtBufferPool.Put(buf)
	buf.Reset()

	p1 := []byte(payload)

	buf.WriteString(hex.EncodeToString(p1))
	buf.WriteByte('.')
	buf.WriteString(hex.EncodeToString(hmac.New(sha256.New, []byte(secret)).Sum(p1)))

	return buf.String()
}

func JWTDecode(secret, token string) ([]byte, error) {
	items := strings.Split(token, ".")
	if len(items) != 2 {
		return nil, fmt.Errorf("not a jwt token")
	}

	payload, err := hex.DecodeString(items[0])
	if err != nil {
		return nil, fmt.Errorf("cannot decode the jwt token")
	}

	signature := hex.EncodeToString(hmac.New(sha256.New, []byte(secret)).Sum(payload))

	if !strings.EqualFold(signature, items[1]) {
		return nil, fmt.Errorf("signature does not match")
	}

	return payload, nil
}
