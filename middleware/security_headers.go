package middleware

import (
	"net/http"

	"github.com/nskforward/httpx/types"
)

func AddSecurityHeaders(next types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// Strict-Transport-Security tells the browser should remember that a site is only to be accessed using HTTPS
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubdomains; preload")

		// Content-Security-Policy tells the browser which kind if resources should be loaded based on origins
		w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; frame-ancestors 'none'")

		// X-Content-Type-Options header allows you to avoid MIME type sniffing by saying that the MIME types are deliberately configured.
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// X-Frame-Options header can be used to indicate whether a browser should be allowed to render a page in a <frame>, <iframe>, <embed> or <object>.
		w.Header().Set("X-Frame-Options", "DENY")

		// X-XSS-Protection header was a feature of Internet Explorer, Chrome and Safari that stopped pages from loading when they detected reflected cross-site scripting (XSS) attacks.
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy header controls how much referrer information (sent with the Referer header) should be included with requests.
		w.Header().Set("Referrer-Policy", "same-origin")

		return next(w, r)
	}
}
