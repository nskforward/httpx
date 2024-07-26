package test

import (
	"fmt"
	"testing"

	"github.com/nskforward/httpx/cache"
)

func TestHash(t *testing.T) {
	s := "/api/v1/users"
	fmt.Println(cache.Hash(s))
	fmt.Println(cache.Hash(s))
}
