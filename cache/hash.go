package cache

import (
	"hash/maphash"
	"strconv"
)

func Hash(s ...string) string {
	var hasher maphash.Hash
	for _, s1 := range s {
		hasher.WriteString(s1)
	}
	return strconv.FormatUint(hasher.Sum64(), 36)
}
