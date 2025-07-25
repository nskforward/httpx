package mux

type Method uint8

const (
	ANY Method = 1 << iota
	GET
	POST
	PUT
	PATCH
	DELETE
	HEAD
	OPTIONS
)

func MethodToUInt8(method string) Method {
	switch method {
	case "ANY":
		return ANY
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "PATCH":
		return PATCH
	case "DELETE":
		return DELETE
	case "HEAD":
		return HEAD
	case "OPTIONS":
		return OPTIONS
	default:
		return 0
	}
}

func MethodToStr(method Method) string {
	switch method {
	case ANY:
		return "ANY"
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case PATCH:
		return "PATCH"
	case DELETE:
		return "DELETE"
	case HEAD:
		return "HEAD"
	case OPTIONS:
		return "OPTIONS"
	default:
		return "UNKNOWN"
	}
}
