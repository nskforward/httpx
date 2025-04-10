package cors

import (
	"net/url"
	"strings"
)

type OriginPool struct {
	origins        []Origin
	allowLocalhost bool
}

func ParseOriginPool(allowLocalhost bool, origins []string) OriginPool {
	result := make([]Origin, 0, len(origins))
	for _, s := range origins {
		origin, err := ParseOrigin(s)
		if err != nil {
			panic(err)
		}
		result = append(result, origin)
	}
	return OriginPool{
		origins:        result,
		allowLocalhost: allowLocalhost,
	}
}

func (p OriginPool) Valid(origin string) bool {
	input, err := url.Parse(origin)
	if err != nil {
		return false
	}
	for _, item := range p.origins {
		if item.Valid(input) {
			return true
		}
	}

	if p.allowLocalhost {
		if strings.HasSuffix(origin, "://localhost") || strings.Contains(origin, "://localhost:") {
			return true
		}
	}

	return false
}
