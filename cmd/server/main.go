package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service/config"
	"service/logging"
	serv "service/server"
	"service/tracing"
)

const (
	// version     = "v1.0"
	serviceName = "template"
)

type CloseFunc func() error

func (f CloseFunc) Close() error {
	err := f()
	if err != nil {
		return err
	}

	return nil
}

type Closer struct {
	logger   *zerolog.Logger
	forClose []io.Closer
}

func NewCloser(logger *zerolog.Logger, forClose ...io.Closer) *Closer {
	return &Closer{logger: logger, forClose: forClose}
}

func (c *Closer) Append(forClose ...io.Closer) {
	c.forClose = append(c.forClose, forClose...)
}

func (c *Closer) Close() {
	for _, closer := range c.forClose {
		err := closer.Close()
		if err != nil {
			c.logger.Err(err).Msg("failed to close:")
		}
	}

	c.logger.Info().Msg("all dependencies are closed")
}

func main() {
	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	c := &config.Config{}
	// config
	err := config.ParseConfig(c)
	if err != nil {
		log.Fatal(err)
	}
	// logger
	logger := logging.New(
		logging.WithTimestamp(),
		logging.WithServiceName(serviceName),
		logging.WithPID(),
	)

	logger.Info().Any("config", c).Send()

	forClose := NewCloser(logger)
	defer forClose.Close()

	if err = tracing.RegisterTracerProvider(ctx, serviceName, c.OTEL.TraceEndpoint); err != nil {
		logger.Fatal().Err(err).Msg("failed to register tracer provider")
	}

	// main server
	mainServiceServer, mainRouter := serv.NewServer(c.Server)

	// register routes
	//		main
	fc := RegisterMainServiceRoutes(mainRouter)

	forClose.Append(fc...)

	_, errCh := serv.Run(ctx, mainServiceServer)

	for err = range errCh {
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Err(err).Send()
		}
	}
}

func Middlewares(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}

func RegisterMainServiceRoutes(r chi.Router) []io.Closer { //nolint:unparam
	// middlewares
	Middlewares(r)

	return nil
}
