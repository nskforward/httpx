package hsts

import (
	"strconv"
	"strings"
	"time"
)

type Config struct {
	MaxAge     time.Duration
	SubDomains bool
}

func (cfg Config) Encode() string {
	arr := make([]string, 0, 3)
	arr = append(arr, strconv.Itoa(int(cfg.MaxAge.Seconds())))
	if cfg.SubDomains {
		arr = append(arr, "includeSubDomains")
	}
	return strings.Join(arr, "; ")
}
