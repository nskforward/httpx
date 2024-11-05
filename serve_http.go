package httpx

import (
	"fmt"
	"net/http"
)

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if router.mux == nil {
		panic(fmt.Errorf("uninitialized router"))
	}
	h := func(ww http.ResponseWriter, rr *http.Request) error {
		router.mux.ServeHTTP(ww, rr)
		return nil
	}
	for i := len(router.middlewares) - 1; i >= 0; i-- {
		h = router.middlewares[i](h)
	}
	router.Catch(h)(w, r)
}
