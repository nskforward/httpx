package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/nskforward/httpx"
)

type CorsOptions struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           time.Duration
	allowAnyOrigin   bool
	maxAgeString     string
}

func CorsMiddleware(options CorsOptions) httpx.Handler {
	for _, origin := range options.AllowOrigins {
		if origin == "*" && len(options.AllowOrigins) > 1 {
			panic("cors: AllowOrigins cannot contain several values with wildcard")
		}
		if origin == "*" {
			continue
		}
		if origin == "null" {
			continue
		}
		if origin == "" {
			panic("cors: AllowOrigins contains empty value")
		}
		if !strings.HasPrefix(origin, "http") {
			panic("cors: AllowOrigins value must begin with 'http'")
		}
		if strings.HasSuffix(origin, "/") {
			panic("cors: AllowOrigins value cannot end with '/'")
		}
	}

	options.maxAgeString = "3600"
	if options.MaxAge > 0 {
		options.maxAgeString = strconv.FormatFloat(options.MaxAge.Seconds(), 'f', 0, 64)
	}

	if len(options.AllowOrigins) == 0 || options.AllowOrigins[0] == "*" {
		options.allowAnyOrigin = true
	}

	if options.allowAnyOrigin && options.AllowCredentials {
		panic("cors: AllowCredentials can be used only with exact origin(s)")
	}

	return func(req *http.Request, resp *httpx.Response) error {
		if req.Method == http.MethodOptions {
			origin := req.Header.Get("Origin")
			if origin == "" {
				return resp.Next(req)
			}
			resp.SetHeader("Vary", "Origin")
			err := corsSendHeaders(resp, options, origin)
			if err != nil {
				return err
			}
			return resp.NoContent()
		}

		if req.Header.Get("Sec-Fetch-Mode") == "cors" {
			origin := req.Header.Get("Origin")
			if origin == "" {
				return resp.Next(req)
			}
			err := corsSendHeaders(resp, options, origin)
			if err != nil {
				return err
			}
		}

		return resp.Next(req)
	}
}

func corsSendHeaders(resp *httpx.Response, options CorsOptions, origin string) error {
	if options.allowAnyOrigin {
		resp.SetHeader("Access-Control-Allow-Origin", "*")
	}

	if !options.allowAnyOrigin && !slices.Contains(options.AllowOrigins, origin) {
		return resp.Text(http.StatusForbidden, fmt.Sprintf("unknown origin: %s", origin))
	}

	resp.SetHeader("Access-Control-Allow-Origin", origin)

	if len(options.AllowMethods) > 0 {
		resp.SetHeader("Access-Control-Allow-Methods", strings.Join(options.AllowMethods, ", "))
	} else {
		resp.SetHeader("Access-Control-Allow-Methods", "*")
	}

	if len(options.AllowedHeaders) > 0 {
		resp.SetHeader("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ", "))
	} else {
		resp.SetHeader("Access-Control-Allow-Headers", "*")
	}

	if options.AllowCredentials {
		resp.SetHeader("Access-Control-Allow-Credentials", "true")
	}

	if len(options.ExposedHeaders) > 0 {
		resp.SetHeader("Access-Control-Expose-Headers", strings.Join(options.ExposedHeaders, ", "))
	} else {
		resp.SetHeader("Access-Control-Expose-Headers", "*")
	}

	resp.SetHeader("Access-Control-Max-Age", options.maxAgeString)
	return nil
}
