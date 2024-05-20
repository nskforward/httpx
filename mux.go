package httpx

import (
	"fmt"
	"net/http"

	"github.com/nskforward/httpx/radix"
)

type ServeMux struct {
	tree  *radix.Node
	on404 http.HandlerFunc
	on405 http.HandlerFunc
}

type muxHandler struct {
	get     http.HandlerFunc
	post    http.HandlerFunc
	put     http.HandlerFunc
	patch   http.HandlerFunc
	delete  http.HandlerFunc
	head    http.HandlerFunc
	options http.HandlerFunc
	any     http.HandlerFunc
}

func (mux *ServeMux) HandlerOn404(h http.HandlerFunc) {
	mux.on404 = h
}

func (mux *ServeMux) HandlerOn405(h http.HandlerFunc) {
	mux.on405 = h
}

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if mux.tree == nil {
		mux.serveError(404, w, r)
		return
	}
	node := mux.tree.Search(path)
	if node == nil {
		mux.serveError(404, w, r)
		return
	}
	mh, ok := node.Value().(*muxHandler)
	if !ok {
		mux.serveError(404, w, r)
		return
	}
	h := mh.lookup(r.Method)
	if h == nil {
		mux.serveError(405, w, r)
		return
	}
	h(w, r)
}

func (mux *ServeMux) serveError(status int, w http.ResponseWriter, r *http.Request) {
	if status == 404 && mux.on404 != nil {
		mux.on404(w, r)
		return
	}

	if status == 405 && mux.on405 != nil {
		mux.on405(w, r)
		return
	}

	http.Error(w, http.StatusText(status), status)
}

func (mux *ServeMux) Route(method, pattern string, h http.HandlerFunc) error {
	if mux.tree == nil {
		mux.tree = &radix.Node{}
	}
	node := mux.tree.Search(pattern)
	if node != nil {
		mh, ok := node.Value().(*muxHandler)
		if ok {
			return mh.registry(method, h)
		}
	}
	mh := &muxHandler{}
	err := mh.registry(method, h)
	if err != nil {
		return err
	}
	return mux.tree.Insert(pattern, mh)
}

func (mh *muxHandler) registry(method string, h http.HandlerFunc) error {
	switch method {
	case "GET":
		if mh.get != nil {
			return fmt.Errorf("method %s already defined", method)
		}
		mh.get = h
		return nil

	case "POST":
		if mh.post != nil {
			return fmt.Errorf("method %s already defined", method)
		}
		mh.post = h
		return nil

	case "PUT":
		if mh.put != nil {
			return fmt.Errorf("method %s already defined", method)
		}
		mh.put = h
		return nil

	case "PATCH":
		if mh.patch != nil {
			return fmt.Errorf("method %s already defined", method)
		}
		mh.patch = h
		return nil

	case "DELETE":
		if mh.delete != nil {
			return fmt.Errorf("method %s already defined", method)
		}
		mh.delete = h
		return nil

	case "HEAD":
		if mh.head != nil {
			return fmt.Errorf("method %s already defined", method)
		}
		mh.head = h
		return nil

	case "OPTIONS":
		if mh.options != nil {
			return fmt.Errorf("method %s already defined", method)
		}
		mh.options = h
		return nil

	case "ANY":
		if mh.any != nil {
			return fmt.Errorf("method %s already defined", method)
		}
		mh.any = h
		return nil

	default:
		return fmt.Errorf("unknown method '%s'", method)
	}
}

func (mh *muxHandler) lookup(method string) http.HandlerFunc {
	switch method {
	case "GET":
		if mh.get != nil {
			return mh.get
		}
		return mh.any

	case "POST":
		if mh.post != nil {
			return mh.post
		}
		return mh.any

	case "PUT":
		if mh.put != nil {
			return mh.put
		}
		return mh.any

	case "PATCH":
		if mh.patch != nil {
			return mh.patch
		}
		return mh.any

	case "DELETE":
		if mh.delete != nil {
			return mh.delete
		}
		return mh.any

	case "HEAD":
		if mh.head != nil {
			return mh.head
		}
		return mh.any

	case "OPTIONS":
		if mh.options != nil {
			return mh.options
		}
		return mh.any

	default:
		return nil
	}
}
