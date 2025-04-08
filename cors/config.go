package cors

import "time"

type Config struct {

	// AllowLocalhost allows localhost origin with any port and schema even AllowOrigins does not contain that origin
	AllowLocalhost bool

	// AllowCredentials allows cookies, client TLS certificate, headers with password
	AllowCredentials bool

	// AllowOrigins allow the origin list for CORS requests
	AllowOrigins []string

	// AllowMethods allows the http method list for CORS requests
	// GET, HEAD, and POST are always allowed
	AllowMethods []string

	// AllowHeaders allows the http header list for CORS requests
	// Accept, Accept-Language, Content-Language, Content-Type, Range are always allowed
	AllowHeaders []string

	// ExposeHeaders allows clients to access the listed headers via javascript
	// Cache-Control, Content-Language, Content-Length, Content-Type, Expires, Last-Modified, Pragma are always allowed
	ExposeHeaders []string

	// MaxAge allows to cache the response of preflight requests for particular period of time
	MaxAge time.Duration
}
