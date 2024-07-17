package transport

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Listener struct {
	ctx              context.Context
	cancel           context.CancelFunc
	l                net.Listener
	maxConnections   int
	maxRateConn      int
	maxRateConnBurst int
	counters         map[string]*Counter
	mx               sync.Mutex
}

func NewListener(ctx context.Context, l net.Listener, maxConnections, maxRateConn int) net.Listener {
	ctx1, cancel := context.WithCancel(ctx)
	return &Listener{
		l:                l,
		maxConnections:   maxConnections,
		maxRateConn:      maxRateConn,
		maxRateConnBurst: maxRateConn * 2,
		ctx:              ctx1,
		cancel:           cancel,
	}
}

func (l *Listener) Close() error {
	l.cancel()
	return l.l.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.l.Addr()
}

func (l *Listener) Accept() (net.Conn, error) {
	c1, err := l.l.Accept()
	if err != nil {
		return nil, err
	}

	ip, _, err := net.SplitHostPort(c1.RemoteAddr().String())
	if err != nil {
		c1.Close()
		return nil, err
	}

	l.mx.Lock()
	counter, ok := l.counters[ip]
	if !ok {
		counter = &Counter{
			limiter: rate.NewLimiter(rate.Every(time.Second/time.Duration(l.maxRateConn)), l.maxRateConnBurst),
		}
		l.counters[ip] = counter
	}
	l.mx.Unlock()
	count := counter.Inc()
	if count > int64(l.maxConnections) {
		c1.Close()
		return nil, fmt.Errorf("reached the max tcp connections (%d) for client: %s", l.maxConnections, ip)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = counter.limiter.Wait(ctx)
	if err != nil {
		c1.Close()
		return nil, fmt.Errorf("connection rate too high (%d) for client: %s", l.maxRateConn, ip)
	}

	return &Conn{Conn: c1, IP: ip, Listener: l}, nil
}
