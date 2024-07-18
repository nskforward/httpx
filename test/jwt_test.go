package test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/nskforward/httpx/jwt"
)

/*
sha256	75736572313031406578616d706c652e636f6d.75736572313031406578616d706c652e636f6dbdc8b3aedc55e68fa073263e195affd415f6e6aa7754fc5f931ca212efc31e6e
sha1	75736572313031406578616d706c652e636f6d.75736572313031406578616d706c652e636f6dff4595f9360d2b7f8fe44bd64702c6d2da11f7d8
*/

func TestJWT(t *testing.T) {
	encoder := jwt.NewEncoder("foobarba")

	email1 := []byte("user1@example.com")
	email2 := []byte("user2@example.com")

	token := encoder.Encode(email1)
	fmt.Println(token)
	payload, err := encoder.Decode([]byte(token))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(payload))

	token = encoder.Encode(email2)
	if err != nil {
		t.Fatal(err)
	}
	payload, err = encoder.Decode([]byte(token))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(payload))
}

/*
BenchmarkJWT-4	182428	      5646 ns/op	    2187 B/op	      21 allocs/op
BenchmarkJWT-4	138412	      7367 ns/op	    3431 B/op	      29 allocs/op
BenchmarkJWT-4	139058	      7313 ns/op	    3625 B/op	      28 allocs/op
BenchmarkJWT-4	215156	      5537 ns/op	    2412 B/op	      22 allocs/op
BenchmarkJWT-4	187386	      5518 ns/op	    2157 B/op	      20 allocs/op
BenchmarkJWT-4	179772	      5747 ns/op	    1842 B/op	      17 allocs/op
BenchmarkJWT-4	168354	      6029 ns/op	    1725 B/op	      15 allocs/op
BenchmarkJWT-4	215650	      4756 ns/op	    1225 B/op	      12 allocs/op
*/
func BenchmarkJWT(b *testing.B) {
	encoder := jwt.NewEncoder("foobarba")
	email := []byte("user101foobarbazz202@example.com")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token := encoder.Encode(email)
			payload, err := encoder.Decode([]byte(token))
			if err != nil {
				b.Fatal(err)
			}
			if !bytes.Equal(email, payload) {
				b.Fatal(fmt.Errorf("not equal"))
			}
		}
	})
}
