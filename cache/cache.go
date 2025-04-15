package cache

import (
	"fmt"
	"time"

	"github.com/nskforward/httpx"
)

func Disallow(resp *httpx.Response) {
	resp.SetHeader("Cache-Control", "no-store")
}

func Allow(resp *httpx.Response, maxAge time.Duration, public bool) {
	scope := func(isPublic bool) string {
		if isPublic {
			return "public"
		}
		return "private"
	}
	resp.SetHeader("Cache-Control", fmt.Sprintf("%s, max-age=%.0f", scope(public), maxAge.Seconds()))
}
