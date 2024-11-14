package jwt

import (
	"time"
)

type Claims struct {
	ID        string
	Subject   string
	Issuer    string
	Audience  []string
	IssuedAt  time.Time
	ExpiresAt time.Time
	NotBefore time.Time
}
