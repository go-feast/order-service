package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/traceid"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service/config"
	"service/http/metrics"
	mw "service/http/middleware"
	httptracing "service/http/tracing"
	"service/logging"
	"service/metrics"
	serv "service/server"
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
	logger   *logging.Logger
	forClose []io.Closer
}

func NewCloser(logger *logging.Logger, forClose ...io.Closer) *Closer {
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

	shutdown, err := httptracing.Provider(context.Background(), serviceName, c.OTEL.TraceEndpoint)
	if err != nil {
		logger.Panic().Err(err).Msg("trace provider error")
	}

	forClose.Append(CloseFunc(func() error { return shutdown(context.Background()) }))

	// metric server
	metricServer, metricRouter := serv.NewServer(c.MetricServer)

	httpmetrics.RegisterServer(
		metrics.NewMetricCollector(metrics.NewMetrics(serviceName), prometheus.NewRegistry()),
		logger,
	)

	// main server
	mainServiceServer, mainRouter := serv.NewServer(c.Server)

	// register routes
	//		main
	fc := RegisterMainServiceRoutes(mainRouter)

	forClose.Append(fc...)

	// 		metrics
	RegisterMetricRoutes(metricRouter)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	_, errCh := serv.Run(ctx, mainServiceServer, metricServer)

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
	r.Use(httpmetrics.RecordRequestHit)
	r.Use(traceid.Middleware)
	r.Use(middleware.Recoverer)

}

func RegisterMainServiceRoutes(r chi.Router) []io.Closer { //nolint:unparam
	// middlewares
	Middlewares(r)

	return nil
}

// RegisterMetricRoutes initialize routes for metrics. Example:
//
// [/metric] - for accessing metrics
//
// [/ping] [/healthz] [/readyz] - for checking if service alive
func RegisterMetricRoutes(r chi.Router) {
	r.Use(middleware.Logger)

	r.Get("/healthz", mw.Healthz)
	r.Get("/readyz", mw.Readyz)
	r.Get("/ping", mw.Ping)

	r.Get("/metrics", httpmetrics.Handler())
}
