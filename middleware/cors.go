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
	allowAnyOrigin   bool
	maxAgeSting      string
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

	options.maxAgeSting = "3600"
	if options.MaxAge > 0 {
		options.maxAgeSting = strconv.FormatFloat(options.MaxAge.Seconds(), 'f', 0, 64)
	}

	if len(options.AllowOrigins) == 0 || options.AllowOrigins[0] == "*" {
		options.allowAnyOrigin = true
	}

	if options.allowAnyOrigin && options.AllowCredentials {
		panic("cors: AllowCredentials can be used only with exact origin(s)")
	}

	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {

			origin := r.Header.Get("Origin")

			if origin == "" {
				w.Header().Set("Vary", "Origin")
				return next(w, r)
			}

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") == "" {
				w.Header().Set("Vary", "Origin")
				return next(w, r)
			}

			err := corsSendHeaders(w, options, origin)
			if err != nil {
				return err
			}

			if r.Method == http.MethodOptions {
				return response.NoContent(w)
			}

			return next(w, r)
		}
	}
}

func corsSendHeaders(w http.ResponseWriter, options CorsOptions, origin string) error {
	if options.allowAnyOrigin {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	if !options.allowAnyOrigin && !slices.Contains(options.AllowOrigins, origin) {
		return response.APIError{Status: http.StatusForbidden, Text: fmt.Sprintf("unknown origin: %s", origin)}
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)

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

	w.Header().Set("Access-Control-Max-Age", options.maxAgeSting)
	return nil
}
