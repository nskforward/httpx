package cors

import (
	"fmt"
	"slices"
	"strings"

	"github.com/nskforward/httpx"
)

func sendAllowOrigin(cfg Config, origin string, resp *httpx.Response) error {

	if slices.Contains(cfg.AllowOrigins, origin) {
		resp.SetHeader("Vary", "Origin")
		resp.SetHeader("Access-Control-Allow-Origin", origin)
		return nil
	}

	if cfg.AllowLocalhost {
		if strings.HasSuffix(origin, "://localhost") || strings.Contains(origin, "://localhost:") {
			resp.SetHeader("Vary", "Origin")
			resp.SetHeader("Access-Control-Allow-Origin", origin)
			return nil
		}
	}

	return fmt.Errorf("cors origin '%s' not allowed", origin)
}

func sendAllowMethods(cfg Config, requestedMethod string, resp *httpx.Response) error {
	if slices.Contains([]string{"GET", "POST", "HEAD"}, requestedMethod) {
		resp.SetHeader("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
		return nil
	}
	if slices.Contains(cfg.AllowMethods, requestedMethod) {
		resp.SetHeader("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
		return nil
	}
	return fmt.Errorf("cors method '%s' not allowed", requestedMethod)
}

func sendAllowHeaders(cfg Config, requestedHeaders string, resp *httpx.Response) error {
	headers := strings.Split(requestedHeaders, ",")

	for _, h := range headers {
		normalized := toTitle(strings.TrimSpace(h))
		if slices.Contains([]string{"Accept", "Accept-Language", "Content-Language", "Content-Type", "Range"}, normalized) {
			continue
		}
		if slices.Contains(cfg.AllowHeaders, normalized) {
			continue
		}
		return fmt.Errorf("cors request header '%s' not allowed", normalized)
	}

	resp.SetHeader("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ", "))
	return nil
}

func sendAllowCredentials(cfg Config, resp *httpx.Response) {
	if cfg.AllowCredentials {
		resp.SetHeader("Access-Control-Allow-Credentials", "true")
	}
}

func sendExposeHeaders(cfg Config, resp *httpx.Response) {
	if len(cfg.ExposeHeaders) > 0 {
		resp.SetHeader("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
	}
}

func sendMaxAge(maxAge string, resp *httpx.Response) {
	if maxAge != "" {
		resp.SetHeader("Access-Control-Max-Age", maxAge)
	}
}

func toTitle(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
