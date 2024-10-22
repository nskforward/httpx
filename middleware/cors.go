package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/nskforward/httpx/response"
	"github.com/nskforward/httpx/types"
)

type CorsOptions struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           time.Duration
}

func Cors(options CorsOptions) types.Middleware {
	if len(options.AllowOrigins) > 0 && options.AllowOrigins[0] == "*" && options.AllowCredentials {
		panic("cors: cannot use wildcard in AllowOrigins with enabled AllowCredentials")
	}

	if len(options.AllowMethods) > 0 && options.AllowMethods[0] == "*" && options.AllowCredentials {
		panic("cors: cannot use wildcard in AllowMethods with enabled AllowCredentials")
	}

	if len(options.AllowedHeaders) > 0 && options.AllowedHeaders[0] == "*" && options.AllowCredentials {
		panic("cors: cannot use wildcard in AllowedHeaders with enabled AllowCredentials")
	}

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

	standardHeaders := []string{"accept", "accept-language", "content-language", "content-type", "range"}
	for _, h := range options.AllowedHeaders {
		if slices.Contains(standardHeaders, strings.ToLower(h)) {
			panic(fmt.Sprintf("cors: AllowedHeaders contains standard header that is always allowed by default: %s", h))
		}
	}

	if len(options.ExposedHeaders) > 0 && options.ExposedHeaders[0] == "*" && options.AllowCredentials {
		panic("cannot use wildcard in Cors.ExposedHeaders with enabled Cors.AllowCredentials")
	}

	maxAge := "3600"
	if options.MaxAge > 0 {
		maxAge = strconv.FormatFloat(options.MaxAge.Seconds(), 'f', 0, 64)
	}

	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {

			origin := r.Header.Get("Origin")
			if origin != "" && !slices.Contains(options.AllowOrigins, origin) {
				return response.APIError{Status: http.StatusForbidden, Text: fmt.Sprintf("unknown origin: %s", origin)}
			}

			if r.Method != http.MethodOptions || r.Header.Get("Access-Control-Request-Method") == "" {
				if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
				return next(w, r)
			}

			if origin == "" {
				return response.APIError{Status: http.StatusForbidden, Text: "origin request header cannot be empty"}
			}

			if !slices.Contains(options.AllowOrigins, origin) {
				return response.APIError{Status: http.StatusForbidden, Text: fmt.Sprintf("unknown origin: %s", origin)}
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			if len(options.AllowMethods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowMethods, ", "))
			}
			if len(options.AllowedHeaders) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ", "))
			}
			if options.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if len(options.ExposedHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(options.ExposedHeaders, ", "))
			}
			w.Header().Set("Access-Control-Max-Age", maxAge)

			return response.NoContent(w)
		}
	}
}
