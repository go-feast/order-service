package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"service/config"
	"sync"
)

func app(ctx context.Context, c *config.ServerConfig, l *zap.Logger) error {
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

func server(c *config.ServerConfig) (server *http.Server, r chi.Router) {
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
