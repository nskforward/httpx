package middleware

import (
	"net/http"

	"github.com/nskforward/httpx/types"
)

type SecurityHeaders struct {
	// Strict-Transport-Security tells the browser should remember that a site is only to be accessed using HTTPS
	StrictTransportSecurity string

	// Content-Security-Policy tells the browser which kind if resources should be loaded based on origins
	ContentSecurityPolicy string

	// X-Content-Type-Options header allows you to avoid MIME type sniffing by saying that the MIME types are deliberately configured.
	ContentTypeOptions string

	// X-Frame-Options header can be used to indicate whether a browser should be allowed to render a page in a <frame>, <iframe>, <embed> or <object>.
	FrameOptions string

	// X-XSS-Protection header was a feature of Internet Explorer, Chrome and Safari that stopped pages from loading when they detected reflected cross-site scripting (XSS) attacks.
	XSSProtection string

	// Referrer-Policy header controls how much referrer information (sent with the Referer header) should be included with requests.
	ReferrerPolicy string
}

func DefaultSecurityHeaders() SecurityHeaders {
	return SecurityHeaders{
		StrictTransportSecurity: "max-age=63072000; includeSubdomains; preload",
		ContentSecurityPolicy:   "default-src 'none'; img-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; frame-ancestors 'none'",
		ContentTypeOptions:      "nosniff",
		FrameOptions:            "DENY",
		XSSProtection:           "1; mode=block",
		ReferrerPolicy:          "same-origin",
	}
}

func AddSecurityHeaders(headers SecurityHeaders) types.Middleware {
	return func(next types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if headers.StrictTransportSecurity != "" {
				w.Header().Set("Strict-Transport-Security", headers.StrictTransportSecurity)
			}
			if headers.ContentSecurityPolicy != "" {
				w.Header().Set("Content-Security-Policy", headers.ContentSecurityPolicy)
			}
			if headers.ContentTypeOptions != "" {
				w.Header().Set("X-Content-Type-Options", headers.ContentTypeOptions)
			}
			if headers.FrameOptions != "" {
				w.Header().Set("X-Frame-Options", headers.FrameOptions)
			}
			if headers.XSSProtection != "" {
				w.Header().Set("X-XSS-Protection", headers.XSSProtection)
			}
			if headers.ReferrerPolicy != "" {
				w.Header().Set("Referrer-Policy", headers.ReferrerPolicy)
			}

			return next(w, r)
		}
	}
}
