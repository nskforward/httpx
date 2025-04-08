package cors

import (
	"slices"
	"strings"

	"github.com/nskforward/httpx"
)

func sendAllowOrigin(cfg Config, origin string, resp *httpx.Response) {
	if slices.Contains(cfg.AllowOrigins, origin) {
		resp.SetHeader("Access-Control-Allow-Origin", origin)
		return
	}

	if cfg.AllowLocalhost {
		if strings.HasSuffix(origin, "://localhost") || strings.Contains(origin, "://localhost:") {
			resp.SetHeader("Access-Control-Allow-Origin", origin)
			return
		}
	}

	if len(cfg.AllowOrigins) == 0 {
		resp.SetHeader("Access-Control-Allow-Origin", "*")
		return
	}

	resp.SetHeader("Access-Control-Allow-Origin", cfg.AllowOrigins[0])
}

func sendAllowMethods(cfg Config, resp *httpx.Response) {
	if len(cfg.AllowMethods) == 0 {
		return
	}
	resp.SetHeader("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
}

func sendAllowHeaders(cfg Config, resp *httpx.Response) {
	if len(cfg.AllowHeaders) > 0 {
		resp.SetHeader("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ", "))
	}
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
