package httpx

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
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

func CorsMiddleware(options CorsOptions) Middleware {
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

	return func(next Handler) Handler {
		return func(ctx *Context) error {

			if ctx.Method() == http.MethodOptions {
				origin := ctx.GetRequestHeader("Origin")
				if origin == "" {
					return next(ctx)
				}
				ctx.SetResponseHeader("Vary", "Origin")
				err := corsSendHeaders(ctx, options, origin)
				if err != nil {
					return err
				}
				return ctx.RespondNoContent()
			}

			if ctx.GetRequestHeader("Sec-Fetch-Mode") == "cors" {
				origin := ctx.GetRequestHeader("Origin")
				if origin == "" {
					return next(ctx)
				}
				err := corsSendHeaders(ctx, options, origin)
				if err != nil {
					return err
				}
			}

			return next(ctx)
		}
	}
}

func corsSendHeaders(ctx *Context, options CorsOptions, origin string) error {
	if options.allowAnyOrigin {
		ctx.SetResponseHeader("Access-Control-Allow-Origin", "*")
	}

	if !options.allowAnyOrigin && !slices.Contains(options.AllowOrigins, origin) {
		return ctx.RespondText(http.StatusForbidden, fmt.Sprintf("unknown origin: %s", origin))
	}

	ctx.SetResponseHeader("Access-Control-Allow-Origin", origin)

	if len(options.AllowMethods) > 0 {
		ctx.SetResponseHeader("Access-Control-Allow-Methods", strings.Join(options.AllowMethods, ", "))
	} else {
		ctx.SetResponseHeader("Access-Control-Allow-Methods", "*")
	}

	if len(options.AllowedHeaders) > 0 {
		ctx.SetResponseHeader("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ", "))
	} else {
		ctx.SetResponseHeader("Access-Control-Allow-Headers", "*")
	}

	if options.AllowCredentials {
		ctx.SetResponseHeader("Access-Control-Allow-Credentials", "true")
	}

	if len(options.ExposedHeaders) > 0 {
		ctx.SetResponseHeader("Access-Control-Expose-Headers", strings.Join(options.ExposedHeaders, ", "))
	} else {
		ctx.SetResponseHeader("Access-Control-Expose-Headers", "*")
	}

	ctx.SetResponseHeader("Access-Control-Max-Age", options.maxAgeString)
	return nil
}
