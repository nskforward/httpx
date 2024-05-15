package httpx

import (
	"net/http"
)

func GET(handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return MethodNotAllowed()
		}
		return handler(w, r)
	}
}

func POST(handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return MethodNotAllowed()
		}
		return handler(w, r)
	}
}

func PUT(handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPut {
			return MethodNotAllowed()
		}
		return handler(w, r)
	}
}

func PATCH(handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPatch {
			return MethodNotAllowed()
		}
		return handler(w, r)
	}
}

func DELETE(handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodDelete {
			return MethodNotAllowed()
		}
		return handler(w, r)
	}
}
