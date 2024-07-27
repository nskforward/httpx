package cache

type Size uint64

const (
	KB Size = 1024
	MB      = 1024 * KB
	GB      = 1024 * MB
	TB      = 1024 * GB
	PB      = 1024 * TB
)
