package mux

import (
	"math/rand"
	"net/http"
	"testing"
)

// BenchmarkMapStr-10    	45185734	        26.43 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMapStr(b *testing.B) {
	m := make(map[string]http.Handler)
	m["GET"] = http.NotFoundHandler()
	m["POST"] = http.NotFoundHandler()
	m["PUT"] = http.NotFoundHandler()
	m["PATCH"] = http.NotFoundHandler()
	m["DELETE"] = http.NotFoundHandler()
	m["CONNECT"] = http.NotFoundHandler()
	m["OPTIONS"] = http.NotFoundHandler()
	m["HEAD"] = http.NotFoundHandler()
	m["TRACE"] = http.NotFoundHandler()
	var h http.Handler
	for b.Loop() {
		n := rand.Intn(9)
		switch n {
		case 0:
			h = m["GET"]
		case 1:
			h = m["POST"]
		case 2:
			h = m["PUT"]
		case 3:
			h = m["PATCH"]
		case 4:
			h = m["DELETE"]
		case 5:
			h = m["CONNECT"]
		case 6:
			h = m["OPTIONS"]
		case 7:
			h = m["HEAD"]
		case 8:
			h = m["TRACE"]
		}
	}
	_ = h
}

// BenchmarkMapUInt8-10    	45549548	        26.27 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMapUInt8(b *testing.B) {
	m := make(map[uint8]http.Handler)
	m[1] = http.NotFoundHandler()
	m[2] = http.NotFoundHandler()
	m[3] = http.NotFoundHandler()
	m[4] = http.NotFoundHandler()
	m[5] = http.NotFoundHandler()
	m[6] = http.NotFoundHandler()
	m[7] = http.NotFoundHandler()
	m[8] = http.NotFoundHandler()
	m[9] = http.NotFoundHandler()
	var h http.Handler
	for b.Loop() {
		n := rand.Intn(9)
		switch n {
		case 0:
			h = m[1]
		case 1:
			h = m[2]
		case 2:
			h = m[3]
		case 3:
			h = m[4]
		case 4:
			h = m[5]
		case 5:
			h = m[6]
		case 6:
			h = m[7]
		case 7:
			h = m[8]
		case 8:
			h = m[9]
		}
	}
	_ = h
}
