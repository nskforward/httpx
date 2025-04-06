package httpx

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	addr      string
	router    *Router
	tlsConfig *tls.Config
}

func NewApp(addr string, logger *slog.Logger, tlsConfig *tls.Config) *App {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return &App{
		addr:      addr,
		router:    NewRouter(logger),
		tlsConfig: tlsConfig,
	}
}

func (app *App) Run() error {
	app.router.logger.Info("starting server", "addr", app.addr, "tls", app.tlsConfig != nil)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := http.Server{
		Addr:              app.addr,
		Handler:           app.router,
		TLSConfig:         app.tlsConfig,
		ReadTimeout:       0,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    10 * 1024,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	if app.tlsConfig != nil {
		return server.ListenAndServeTLS("", "")
	}
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (app *App) Router() *Router {
	return app.router
}
