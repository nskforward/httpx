package cors

import (
	"strconv"
)

func NormalizeMaxAge(cfg Config) string {
	if cfg.MaxAge > 0 {
		return strconv.Itoa(int(cfg.MaxAge.Seconds()))
	}
	return ""
}
