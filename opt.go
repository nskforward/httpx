package httpx

import (
	"crypto/tls"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type SetOpt func(*App)

func WithAddr(addr string) SetOpt {
	return func(a *App) {
		a.addr = addr
	}
}

func WithLogger(logger *slog.Logger) SetOpt {
	return func(a *App) {
		a.logger = logger
	}
}

func WithTLSConfig(tlsConfig *tls.Config) SetOpt {
	return func(a *App) {
		a.tlsConfig = tlsConfig
	}
}

func WithReadHeaderTimeout(timeout time.Duration) SetOpt {
	return func(a *App) {
		a.readHeaderTimeout = timeout
	}
}

func WithReadTimeout(timeout time.Duration) SetOpt {
	return func(a *App) {
		a.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) SetOpt {
	return func(a *App) {
		a.writeTimeout = timeout
	}
}

func WithKeepAlive(age time.Duration) SetOpt {
	return func(a *App) {
		a.keepAlive = age
	}
}

func WithMaxHeaderBytes(size int) SetOpt {
	return func(a *App) {
		a.maxHeaderBytes = size
	}
}

func WithConnState(handler func(net.Conn, http.ConnState)) SetOpt {
	return func(a *App) {
		a.connState = handler
	}
}
