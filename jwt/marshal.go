package jwt

import (
	"bytes"
	"strconv"
	"time"
)

func (claims Claims) MarshalJSON() (b []byte, err error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	separator := false
	if claims.ID != "" {
		jsonWriteSeparator(&buf, &separator)
		jsonWriteFieldName(&buf, "jti")
		jsonWriteValueString(&buf, claims.ID)
	}
	if claims.Subject != "" {
		jsonWriteSeparator(&buf, &separator)
		jsonWriteFieldName(&buf, "sub")
		jsonWriteValueString(&buf, claims.Subject)
	}
	if claims.Issuer != "" {
		jsonWriteSeparator(&buf, &separator)
		jsonWriteFieldName(&buf, "iss")
		jsonWriteValueString(&buf, claims.Issuer)
	}
	if len(claims.Audience) > 0 {
		jsonWriteSeparator(&buf, &separator)
		jsonWriteFieldName(&buf, "aud")
		jsonWriteValueStringSlice(&buf, claims.Audience)
	}
	if !claims.IssuedAt.IsZero() {
		jsonWriteSeparator(&buf, &separator)
		jsonWriteFieldName(&buf, "iat")
		jsonWriteValueTime(&buf, claims.IssuedAt)
	}
	if !claims.ExpiresAt.IsZero() {
		jsonWriteSeparator(&buf, &separator)
		jsonWriteFieldName(&buf, "exp")
		jsonWriteValueTime(&buf, claims.ExpiresAt)
	}
	if !claims.NotBefore.IsZero() {
		jsonWriteSeparator(&buf, &separator)
		jsonWriteFieldName(&buf, "nbf")
		jsonWriteValueTime(&buf, claims.NotBefore)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func jsonWriteSeparator(buf *bytes.Buffer, separator *bool) {
	if *separator {
		buf.WriteString(",")
	} else {
		*separator = true
	}
}

func jsonWriteFieldName(buf *bytes.Buffer, field string) {
	buf.WriteString("\"")
	buf.WriteString(field)
	buf.WriteString("\":")
}

func jsonWriteValueString(buf *bytes.Buffer, value string) {
	var prev rune
	buf.WriteRune('"')
	for _, c := range value {
		if c == '"' && prev != '\\' {
			buf.WriteRune('\\')
		}
		buf.WriteRune(c)
		prev = c
	}
	buf.WriteRune('"')
}

func jsonWriteValueStringSlice(buf *bytes.Buffer, values []string) {
	buf.WriteRune('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteRune(',')
		}
		jsonWriteValueString(buf, v)
	}
	buf.WriteRune(']')
}

func jsonWriteValueTime(buf *bytes.Buffer, value time.Time) {
	buf.WriteString(strconv.FormatInt(value.Unix(), 10))
}
