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
			if origin != "" {
				if len(options.AllowOrigins) > 0 {
					if !slices.Contains(options.AllowOrigins, origin) {
						return response.APIError{Status: http.StatusForbidden, Text: fmt.Sprintf("unknown origin: %s", origin)}
					}
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}

				if len(options.AllowMethods) > 0 {
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowMethods, ", "))
				} else {
					w.Header().Set("Access-Control-Allow-Methods", "*")
				}

				if len(options.AllowedHeaders) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ", "))
				} else {
					w.Header().Set("Access-Control-Allow-Headers", "*")
				}

				if options.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				if len(options.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(options.ExposedHeaders, ", "))
				} else {
					w.Header().Set("Access-Control-Expose-Headers", "*")
				}

				w.Header().Set("Access-Control-Max-Age", maxAge)
			}

			if r.Method != http.MethodOptions || r.Header.Get("Access-Control-Request-Method") == "" {
				return next(w, r)
			}

			if origin == "" {
				return response.APIError{Status: http.StatusForbidden, Text: "origin request header cannot be empty"}
			}
			return response.NoContent(w)
		}
	}
}
