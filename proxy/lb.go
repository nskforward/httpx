package proxy

import (
	"math/rand"
	"net/http"

	"github.com/nskforward/httpx/types"
)

func LoadBalancer(backends []string) types.Handler {
	if len(backends) == 0 {
		panic("load balancer must have at least one backend")
	}

	alive := make([]types.Handler, 0, len(backends))
	for _, backend := range backends {
		alive = append(alive, Reverse(backend))
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		handler := alive[rand.Intn(len(alive))]
		return handler(w, r)
	}
}
