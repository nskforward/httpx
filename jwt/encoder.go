package jwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"net/http"
	"slices"
	"sync"
	"unsafe"

	"github.com/nskforward/httpx/types"
)

type Encoder struct {
	buffers sync.Pool
	b       sync.Pool
	hasher  sync.Pool
}

func NewEncoder(secret string) *Encoder {
	return &Encoder{
		buffers: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
		b: sync.Pool{
			New: func() any {
				b := make([]byte, 256)
				return &b
			},
		},
		hasher: sync.Pool{
			New: func() any {
				return hmac.New(sha256.New, []byte(secret))
			},
		},
	}
}

func (enc *Encoder) ParseRequest(r *http.Request) (string, error) {
	token := r.Header.Get(types.Authorization)
	if token == "" {
		return "", fmt.Errorf("require Authorization header")
	}
	data, err := enc.Decode([]byte(token))
	if err != nil {
		return "", fmt.Errorf("bad Authorization header")
	}
	if len(data) == 0 {
		return "", fmt.Errorf("empty decoded value")
	}
	return string(data), nil
}

func (enc *Encoder) Encode(src []byte) string {
	buf := enc.buffers.Get().(*bytes.Buffer)
	buf.Reset()
	defer enc.buffers.Put(buf)

	hasher := enc.hasher.Get().(hash.Hash)
	signature := hasher.Sum(src)
	enc.hasher.Put(hasher)

	bPtr := enc.b.Get().(*[]byte)

	bPtr = hexEnc(bPtr, src)
	buf.Write(*bPtr)
	buf.WriteByte('.')
	bPtr = hexEnc(bPtr, signature)
	buf.Write(*bPtr)

	enc.b.Put(bPtr)

	return buf.String()
}

func (enc *Encoder) Decode(src []byte) ([]byte, error) {
	index := bytes.Index(src, []byte("."))
	if index < 0 {
		return nil, fmt.Errorf("incorrect format")
	}

	encodedPayload := src[:index]

	payload, err := hex.DecodeString(b2s(encodedPayload))
	if err != nil {
		return nil, fmt.Errorf("incorrect format")
	}

	hasher := enc.hasher.Get().(hash.Hash)
	signature := hasher.Sum(payload)
	enc.hasher.Put(hasher)

	bPtr := enc.b.Get().(*[]byte)
	defer enc.b.Put(bPtr)

	bPtr = hexEnc(bPtr, signature)

	if !bytes.Equal(*bPtr, src[index+1:]) {
		return nil, fmt.Errorf("invalid signature")
	}

	return payload, nil
}

func hexEnc(bPtr *[]byte, src []byte) *[]byte {
	dst := *bPtr
	minCap := hex.EncodedLen(len(src))
	if cap(dst) < minCap {
		dst = slices.Grow(dst, minCap)
	}
	dst = dst[:minCap]
	hex.Encode(dst, src)
	return &dst
}

func s2b(str string) []byte {
	if str == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func b2s(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
