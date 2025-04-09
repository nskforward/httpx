package cors

import (
	"fmt"
	"net/url"
	"strings"
)

// *.example.com

func normalizeOrigins(origins []string) []string {
	for i, s := range origins {
		if len(s) == 0 {
			panic("cors origin list contains empty item")
		}
		s = strings.ToLower(strings.TrimSpace(s))
		if s[len(s)-1] == '*' {
			panic("cors origins with trailing wildcard not allowed")
		}
		if s[0] == '*' {
			if strings.Count(s, ".") < 2 {
				panic("cors origins with leading wildcard must be the third level domain or higher for security reason")
			}
			origins[i] = s
			continue
		}
		u, err := url.Parse(s)
		if err != nil {
			panic(fmt.Errorf("cors orign bad format: %w", err))
		}
		if u.Scheme == "" || u.Path != "" {
			panic("cors origin bad format for value '%s', you should follow format: https://example.com, https://example.com:8080")
		}
		origins[i] = s
	}
	return origins
}
