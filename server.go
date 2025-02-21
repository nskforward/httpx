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

type Server struct {
	addr              string
	logger            *slog.Logger
	tlsConfig         *tls.Config
	router            *router
	readHeaderTimeout time.Duration // read request headers
	readTimeout       time.Duration // read whole request
	writeTimeout      time.Duration // write whole response
	keepAlive         time.Duration
	maxHeaderBytes    int // from request
	connState         func(net.Conn, http.ConnState)
}

type SetOpt func(*Server)

func NewServer(opts ...SetOpt) *Server {
	s := &Server{
		addr:      ":80",
		tlsConfig: nil,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.logger == nil {
		s.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	s.router = newRouter(s)
	return s
}

func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func (s *Server) Run() error {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer := http.Server{
		Addr:              s.addr,
		Handler:           s.router,
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
	return httpServer.ListenAndServe()
}

func (s *Server) Use(middlewares ...Handler) {
	s.router.use(middlewares)
}

func (s *Server) Route(method Method, pattern string, handler Handler, middlewares ...Handler) *Route {
	return s.router.Route(method, pattern, handler, middlewares)
}

func (s *Server) Group(pattern string, middlewares ...Handler) *Route {
	return s.router.Group(pattern, middlewares)
}

func (s *Server) GET(pattern string, handler Handler, middlewares ...Handler) *Route {
	return s.router.Route(GET, pattern, handler, middlewares)
}

func (s *Server) POST(pattern string, handler Handler, middlewares ...Handler) *Route {
	return s.router.Route(POST, pattern, handler, middlewares)
}

func (s *Server) PUT(pattern string, handler Handler, middlewares ...Handler) *Route {
	return s.router.Route(PUT, pattern, handler, middlewares)
}

func (s *Server) DELETE(pattern string, handler Handler, middlewares ...Handler) *Route {
	return s.router.Route(DELETE, pattern, handler, middlewares)
}

func (s *Server) PATCH(pattern string, handler Handler, middlewares ...Handler) *Route {
	return s.router.Route(PATCH, pattern, handler, middlewares)
}

func (s *Server) OPTIONS(pattern string, handler Handler, middlewares ...Handler) *Route {
	return s.router.Route(OPTIONS, pattern, handler, middlewares)
}
