package cache

import (
	"strconv"
	"strings"
	"time"
)

type Control struct {
	IsPublic bool
	MaxAge   time.Duration
}

func NewControl(s string) Control {
	cc := Control{}

	items := strings.Split(s, ",")
	for _, item := range items {
		v := strings.ToLower(strings.TrimSpace(item))

		if strings.EqualFold("public", v) {
			cc.IsPublic = true
			continue
		}

		if strings.HasPrefix(v, "max-age=") {
			items2 := strings.Split(v, "=")
			if len(items2) != 2 {
				continue
			}

			v2 := strings.TrimSpace(items2[1])
			n, err := strconv.Atoi(v2)
			if err == nil {
				cc.MaxAge = time.Duration(n) * time.Second
			}

			continue
		}

	}

	return cc
}
