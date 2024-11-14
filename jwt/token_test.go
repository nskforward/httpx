package jwt

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	t1 := NewWithClaims(HS256, Claims{
		Subject: "1234567890",
	})
	out, err := t1.Sign("12345678")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(out)

	t2, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}

	err = t2.Verify("12345678")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("ok")
}
