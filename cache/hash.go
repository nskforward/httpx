package cache

import (
	"hash/crc64"
	"strconv"
)

/*
func Hash(s ...string) string {
	var hasher maphash.Hash

	for _, s1 := range s {
		hasher.WriteString(s1)
	}
	return strconv.FormatUint(hasher.Sum64(), 16)
}
*/

var table = crc64.MakeTable(crc64.ECMA)

func Hash(s ...string) string {
	hasher := crc64.New(table)
	for _, s1 := range s {
		hasher.Write([]byte(s1))
	}
	return strconv.FormatUint(hasher.Sum64(), 16)
}
