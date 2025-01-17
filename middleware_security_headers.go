package httpx

func AddSecurityHeadersMiddleware(next Handler) Handler {
	return func(ctx *Context) error {
		// Strict-Transport-Security tells the browser should remember that a site is only to be accessed using HTTPS
		ctx.SetResponseHeader("Strict-Transport-Security", "max-age=63072000; includeSubdomains; preload")

		// Content-Security-Policy tells the browser which kind if resources should be loaded based on origins
		ctx.SetResponseHeader("Content-Security-Policy", "default-src 'none'; img-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; frame-ancestors 'none'")

		// X-Content-Type-Options header allows you to avoid MIME type sniffing by saying that the MIME types are deliberately configured.
		ctx.SetResponseHeader("X-Content-Type-Options", "nosniff")

		// X-Frame-Options header can be used to indicate whether a browser should be allowed to render a page in a <frame>, <iframe>, <embed> or <object>.
		ctx.SetResponseHeader("X-Frame-Options", "DENY")

		// X-XSS-Protection header was a feature of Internet Explorer, Chrome and Safari that stopped pages from loading when they detected reflected cross-site scripting (XSS) attacks.
		ctx.SetResponseHeader("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy header controls how much referrer information (sent with the Referer header) should be included with requests.
		ctx.SetResponseHeader("Referrer-Policy", "same-origin")

		return next(ctx)
	}
}
