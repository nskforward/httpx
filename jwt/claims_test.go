package jwt

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestClaims(t *testing.T) {

	issuedAt, err := time.Parse("2006-01-02 15:04:05", "2024-11-14 14:34:00")
	if err != nil {
		t.Fatal(err)
	}
	expiresAt, err := time.Parse("2006-01-02 15:04:05", "2025-11-14 14:34:00")
	if err != nil {
		t.Fatal(err)
	}

	claims1 := Claims{
		ID:      "123456789",
		Subject: "user_id",
		Issuer:  "algolego auth service 1.0",
		//Audience:  []string{"read", "write"},
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	data, err := json.Marshal(claims1)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(data))

	var claims2 Claims
	err = json.Unmarshal(data, &claims2)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(claims2)
}
