package httpx

type Method string

const (
	ANY     Method = "ANY"
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
	PATCH   Method = "PATCH"
	HEAD    Method = "HEAD"
	TRACE   Method = "TRACE"
	CONNECT Method = "CONNECT"
)
