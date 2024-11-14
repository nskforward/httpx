package jwt

import (
	"encoding/json"
	"fmt"
	"time"
)

func (claims *Claims) UnmarshalJSON(data []byte) (err error) {
	m := make(map[string]any)
	err = json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	for f, v := range m {
		switch f {

		case "jti":
			claims.ID = v.(string)
			continue

		case "sub":
			claims.Subject = v.(string)
			continue

		case "iss":
			claims.Issuer = v.(string)
			continue

		case "aud":
			slice := v.([]any)
			claims.Audience = make([]string, 0, len(slice))
			for _, s := range slice {
				claims.Audience = append(claims.Audience, s.(string))
			}
			continue

		case "iat":
			claims.IssuedAt = time.Unix(int64(v.(float64)), 0)
			continue

		case "exp":
			claims.ExpiresAt = time.Unix(int64(v.(float64)), 0)
			continue

		case "nbf":
			claims.NotBefore = time.Unix(int64(v.(float64)), 0)
			continue

		default:
			return fmt.Errorf("unknown field: %s", f)
		}
	}

	return nil
}
