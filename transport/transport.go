package transport

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/netutil"
)

type Transport struct {
	// ReadTimeout is max duration to wait the request headers and body before an error. Value must be 0 (disabled) for streaming.
	ReadTimeout time.Duration

	// WriteTimeout is max duration to wait the response headers and body are sent before an error. Value must be 0 (disabled) for streaming.
	WriteTimeout time.Duration

	// ReadHeaderTimeout is max duration to wait the request headers.
	ReadHeaderTimeout time.Duration

	// MaxHeaderBytes is max size of bytes to request headers.
	MaxHeaderBytes int

	// IdleTimeout is duration to wait the next request before connection closed.
	IdleTimeout time.Duration

	// MaxConnectionsTotal is max number of total active connections. Add "LimitNOFILE=100000" Linux setting accordingly for MaxTotalConnections=100000. By default 0 (unlimited).
	MaxConnectionsTotal int

	// MaxConnectionsClient is max number of incomming connections per a period of time for an IP. By defaul unlimited.
	MaxConnectionsClient int

	// MaxConnectionsPerSecondClient limits max connections per second for a client
	MaxConnectionsPerSecondClient int
}

func DefaultTransport() Transport {
	return Transport{
		ReadTimeout:                   0,
		WriteTimeout:                  0,
		ReadHeaderTimeout:             15 * time.Second,
		MaxHeaderBytes:                4096,
		IdleTimeout:                   time.Minute,
		MaxConnectionsTotal:           65536,
		MaxConnectionsClient:          128,
		MaxConnectionsPerSecondClient: 8,
	}
}

func (s Transport) Listen(addr string, h http.Handler) error {
	raw := http.Server{
		Addr:              addr,
		Handler:           h,
		ErrorLog:          log.New(io.Discard, "", 0),
		MaxHeaderBytes:    s.MaxHeaderBytes,
		ReadHeaderTimeout: s.ReadHeaderTimeout,
		IdleTimeout:       s.IdleTimeout,
		ReadTimeout:       s.ReadTimeout,
		WriteTimeout:      s.WriteTimeout,
	}

	ctx, cancel := context.WithCancel(context.Background())

	raw.RegisterOnShutdown(func() {
		cancel()
	})

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if s.MaxConnectionsTotal > 0 {
		l = netutil.LimitListener(l, s.MaxConnectionsTotal)
	}

	l = NewListener(ctx, l, s.MaxConnectionsClient, s.MaxConnectionsPerSecondClient)
	defer l.Close()

	return raw.Serve(l)
}

func (s Transport) ListenTLS(addr string, tlsConfig *tls.Config, h http.Handler) error {
	raw := http.Server{
		Addr:              addr,
		Handler:           h,
		ErrorLog:          log.New(io.Discard, "", 0),
		MaxHeaderBytes:    s.MaxHeaderBytes,
		ReadHeaderTimeout: s.ReadHeaderTimeout,
		IdleTimeout:       s.IdleTimeout,
		ReadTimeout:       s.ReadTimeout,
		WriteTimeout:      s.WriteTimeout,
		TLSConfig:         tlsConfig,
	}

	ctx, cancel := context.WithCancel(context.Background())

	raw.RegisterOnShutdown(func() {
		cancel()
	})

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if s.MaxConnectionsTotal > 0 {
		l = netutil.LimitListener(l, s.MaxConnectionsTotal)
	}

	l = NewListener(ctx, l, s.MaxConnectionsClient, s.MaxConnectionsPerSecondClient)
	defer l.Close()

	return raw.ServeTLS(l, "", "")
}
