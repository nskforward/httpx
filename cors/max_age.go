package cors

import (
	"strconv"
	"time"
)

func NormalizeMaxAge(cfg Config) string {
	if cfg.MaxAge == 0 {
		cfg.MaxAge = time.Minute
	}
	return strconv.Itoa(int(cfg.MaxAge.Seconds()))
}
