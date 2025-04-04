package httpx

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	addr              string
	logger            *slog.Logger
	tlsConfig         *tls.Config
	group             *Group
	readHeaderTimeout time.Duration // read request headers
	readTimeout       time.Duration // read whole request
	writeTimeout      time.Duration // write whole response
	keepAlive         time.Duration
	maxHeaderBytes    int // from request
	connState         func(net.Conn, http.ConnState)
}

func New(opts ...SetOpt) *App {
	s := &App{
		addr:      ":80",
		tlsConfig: nil,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.logger == nil {
		s.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	s.group = &Group{
		r: newRouter(s),
	}
	return s
}

func (s *App) Addr() string {
	return s.addr
}

func (s *App) Handler() http.Handler {
	return s.group.r
}

func (s *App) Run() error {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer := http.Server{
		Addr:              s.addr,
		Handler:           s.group.r,
		TLSConfig:         s.tlsConfig,
		ReadTimeout:       s.readTimeout,
		ReadHeaderTimeout: s.readHeaderTimeout,
		WriteTimeout:      s.writeTimeout,
		IdleTimeout:       s.keepAlive,
		MaxHeaderBytes:    s.maxHeaderBytes,
		BaseContext: func(_ net.Listener) context.Context {
			return mainCtx
		},
		ConnState: s.connState,
	}

	go func() {
		<-mainCtx.Done()
		httpServer.Shutdown(context.Background())
	}()

	if s.tlsConfig != nil {
		return httpServer.ListenAndServeTLS("", "")
	}
	err := httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *App) Use(middlewares ...Handler) {
	s.group.Use(middlewares...)
}

func (s *App) Group(pattern string, middlewares ...Handler) Group {
	return s.group.Group(pattern, middlewares...)
}

func (s *App) Custom(method, pattern string, handler Handler, middlewares ...Handler) {
	s.group.Custom(method, pattern, handler, middlewares...)
}

func (s *App) ANY(pattern string, handler Handler, middlewares ...Handler) {
	s.group.ANY(pattern, handler, middlewares...)
}

func (s *App) GET(pattern string, handler Handler, middlewares ...Handler) {
	s.group.GET(pattern, handler, middlewares...)
}

func (s *App) POST(pattern string, handler Handler, middlewares ...Handler) {
	s.group.POST(pattern, handler, middlewares...)
}

func (s *App) PUT(pattern string, handler Handler, middlewares ...Handler) {
	s.group.PUT(pattern, handler, middlewares...)
}

func (s *App) DELETE(pattern string, handler Handler, middlewares ...Handler) {
	s.group.DELETE(pattern, handler, middlewares...)
}

func (s *App) PATCH(pattern string, handler Handler, middlewares ...Handler) {
	s.group.PATCH(pattern, handler, middlewares...)
}

func (s *App) OPTIONS(pattern string, handler Handler, middlewares ...Handler) {
	s.group.OPTIONS(pattern, handler, middlewares...)
}
