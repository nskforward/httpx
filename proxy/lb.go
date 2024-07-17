package proxy

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/nskforward/httpx/response"
)

type LB struct {
	alive     []*Backend
	dead      []*Backend
	weightSum int
	sync.RWMutex
}

type Backend struct {
	PrefixURL string
	Weight    int
	rp        *httputil.ReverseProxy
}

func NewLB(backends []Backend) *LB {
	if len(backends) == 0 {
		panic("load balancer must have at least one backend")
	}
	lb := &LB{
		alive: make([]*Backend, 0, len(backends)),
	}
	for _, backend := range backends {
		lb.weightSum += backend.Weight
		backend.rp = ReverseProxy(backend.PrefixURL)
		lb.alive = append(lb.alive, &backend)
	}
	if lb.weightSum < 1 {
		panic("backends weight sum must be greater than 0")
	}
	return lb
}

func (lb *LB) ChangeProxyRequest(f func(req *httputil.ProxyRequest)) {
	for _, backend := range lb.alive {
		backend.rp.Rewrite = f
	}
}

func (lb *LB) ChangeProxyResponse(f func(resp *http.Response) error) {
	for _, backend := range lb.alive {
		backend.rp.ModifyResponse = f
	}
}

func (lb *LB) HandleProxyError(f func(w http.ResponseWriter, r *http.Request, err error)) {
	for _, backend := range lb.alive {
		backend.rp.ErrorHandler = f
	}
}

func (lb *LB) HealthCheck(interval time.Duration, urlPath string, header http.Header) {
	total := make([]*Backend, 0, len(lb.alive))
	total = append(total, lb.alive...)

	for _, backend := range total {
		time.Sleep(100 * time.Millisecond)
		go func(backendURL string) {
			target, err := url.JoinPath(backendURL, urlPath)
			if err != nil {
				panic(err)
			}
			ticker := time.NewTicker(interval)
			for range ticker.C {
				err := lb.makeHealthRequest(target, header)
				if err != nil {
					lb.setBackendDead(backendURL)
				} else {
					lb.setBackendAlive(backendURL)
				}
			}
		}(backend.PrefixURL)
	}
}

func (lb *LB) setBackendAlive(backendURL string) {
	backendAlreadyAlive := false
	lb.RLock()
	for _, b1 := range lb.alive {
		if b1.PrefixURL == backendURL {
			backendAlreadyAlive = true
			break
		}
	}
	lb.RUnlock()
	if backendAlreadyAlive {
		return
	}
	var b2 *Backend
	lb.Lock()
	for index, b1 := range lb.dead {
		if b1.PrefixURL == backendURL {
			b2 = lb.dead[index]
			lb.dead[index] = lb.dead[len(lb.dead)-1]
			lb.dead = lb.dead[:len(lb.dead)-1]
			break
		}
	}
	if b2 != nil {
		lb.alive = append(lb.alive, b2)
	}
	lb.weightSum += b2.Weight
	lb.Unlock()
	slog.Info("backend is alive again", "origin", backendURL)
}

func (lb *LB) setBackendDead(backendURL string) {
	backendAlreadyDead := false
	lb.RLock()
	for _, b1 := range lb.dead {
		if b1.PrefixURL == backendURL {
			backendAlreadyDead = true
			break
		}
	}
	lb.RUnlock()
	if backendAlreadyDead {
		return
	}
	var b2 *Backend
	lb.Lock()
	for index, b1 := range lb.alive {
		if b1.PrefixURL == backendURL {
			b2 = lb.alive[index]
			lb.alive[index] = lb.alive[len(lb.alive)-1]
			lb.alive = lb.alive[:len(lb.alive)-1]
			break
		}
	}
	if b2 != nil {
		lb.dead = append(lb.dead, b2)
	}
	lb.weightSum -= b2.Weight
	lb.Unlock()
	slog.Error("backend is dead", "origin", backendURL)
}

func (lb *LB) makeHealthRequest(healthURL string, header http.Header) error {
	req, err := http.NewRequest(http.MethodGet, healthURL, nil)
	if err != nil {
		panic(err)
	}
	for k, vv := range header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, req.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad response code: %s", resp.Status)
	}
	return nil
}

func (lb *LB) roundRobin() (*Backend, error) {
	lb.RWMutex.RLock()
	defer lb.RWMutex.RUnlock()

	random := rand.Intn(lb.weightSum)
	sum := 0
	for _, backend := range lb.alive {
		sum += backend.Weight
		if sum > random {
			return backend, nil
		}
	}
	return nil, response.Error{Status: http.StatusBadGateway, Text: "no alive backends"}
}

func (lb *LB) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	backend, err := lb.roundRobin()
	if err != nil {
		return err
	}
	backend.rp.ServeHTTP(w, r)
	return nil
}
