package cors

import (
	"fmt"
	"slices"
	"strings"

	"github.com/nskforward/httpx"
)

func sendAllowOrigin(cfg Config, origin string, resp *httpx.Response) error {
	if slices.Contains(cfg.AllowOrigins, origin) {
		resp.SetHeader("Access-Control-Allow-Origin", origin)
		return nil
	}

	if cfg.AllowLocalhost {
		if strings.HasSuffix(origin, "://localhost") || strings.Contains(origin, "://localhost:") {
			resp.SetHeader("Access-Control-Allow-Origin", origin)
			return nil
		}
	}

	if len(cfg.AllowOrigins) == 0 && !cfg.AllowCredentials {
		resp.SetHeader("Access-Control-Allow-Origin", "*")
		return nil
	}

	return fmt.Errorf("cors origin '%s' not allowed", origin)
}

func sendAllowMethods(cfg Config, requestedMethod string, resp *httpx.Response) error {
	if len(cfg.AllowMethods) == 0 && !cfg.AllowCredentials {
		resp.SetHeader("Access-Control-Allow-Methods", "*")
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

	if len(cfg.AllowHeaders) == 0 && !cfg.AllowCredentials {
		resp.SetHeader("Access-Control-Allow-Headers", "*")
		return nil
	}
	for _, h := range headers {
		if !slices.Contains(cfg.AllowHeaders, strings.TrimSpace(h)) {
			return fmt.Errorf("cors request header '%s' not allowed", h)
		}
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
	resp.SetHeader("Access-Control-Max-Age", maxAge)
	resp.SetHeader("Vary", "Origin")
}
