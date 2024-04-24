package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"service/config"
	mw "service/http/middleware"
	"service/logging"
	"sync"
)

func main() {
	c := &config.Config{}
	// config
	err := config.ParseEnvironment(c)
	if err != nil {
		log.Fatal(err)
	}
	// logger
	logger, err := logging.NewLogger(c.Environment,
		logging.WithOutputPaths("stdout"),
		logging.WithErrorOutputPaths("stderr"),
	)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("config", zap.Any("config", c))

	// app
	ctx := context.Background()
	if err = app(ctx, c, logger); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("app", zap.Error(err))
	}
}

func app(ctx context.Context, c *config.Config, l *zap.Logger) error {
	// tracing
	// metrics
	// deps: db conn, kafka/rabbit conn
	// server
	l.Info("initializing server")

	s, r := server(c)
	// routes
	l.Info("initializing routes and middlewares")

	closers := decorators(r,
		addMiddleware,
		addRoutes,
	)
	// graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("server shutdown", zap.Error(err))
			return
		}

		l.Info("server shutting down...")
	}()

	// blocking
	l.Info("server is running", zap.String("url", "http://"+net.JoinHostPort(c.Host, c.Port)))

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// wait for Shutdown function to finish its work
	wg.Wait()

	l.Info("closing dependencies")

	var err error

	// safe close deps
	for _, closer := range closers {
		err = closer.Close()
		if err != nil {
			err = errors.Join(err)
		}
	}

	return err
}

func server(c *config.Config) (server *http.Server, r chi.Router) {
	r = chi.NewRouter()

	server = &http.Server{
		Addr:         net.JoinHostPort(c.Host, c.Port),
		Handler:      r,
		WriteTimeout: c.WriteTimeout,
		ReadTimeout:  c.ReadTimeout,
		IdleTimeout:  c.IdleTimeout,
	}

	return
}

func decorators(r chi.Router, dec ...Decorator) []io.Closer {
	arr := make([]io.Closer, 0)

	for _, decorator := range dec {
		closers := decorator(r)
		if len(closers) == 0 {
			continue
		}

		arr = append(arr, closers...)
	}

	return arr
}

type Decorator func(router chi.Router) []io.Closer

func addMiddleware(r chi.Router) []io.Closer {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return nil
}

func healthEndpoints(r chi.Router) []io.Closer { //nolint:unparam
	r.Get("/healthz", mw.Healthz)
	r.Get("/readyz", mw.Readyz)
	r.Get("/ping", mw.Ping)

	return nil
}

func addRoutes(r chi.Router) []io.Closer {
	healthEndpoints(r)
	return nil
}
