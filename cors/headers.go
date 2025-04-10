package cors

import (
	"fmt"
	"slices"
	"strings"

	"github.com/nskforward/httpx"
)

func sendAllowOrigin(resp *httpx.Response, origin string) {
	resp.SetHeader("Vary", "Origin")
	resp.SetHeader("Access-Control-Allow-Origin", origin)
}

func sendAllowMethods(cfg Config, requestedMethod string, resp *httpx.Response) error {
	if slices.ContainsFunc(cfg.AllowMethods, func(s string) bool {
		return strings.EqualFold(s, requestedMethod)
	}) {
		resp.SetHeader("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
		return nil
	}
	return fmt.Errorf("cors method '%s' not allowed", requestedMethod)
}

func sendAllowHeaders(cfg Config, requestedHeaders string, resp *httpx.Response) error {
	headers := strings.Split(requestedHeaders, ",")
	for _, h := range headers {
		normalized := strings.TrimSpace(h)
		if slices.ContainsFunc(cfg.AllowHeaders, func(s string) bool {
			return strings.EqualFold(s, normalized)
		}) {
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
