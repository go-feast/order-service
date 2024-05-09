package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service/config"
	"service/http/metrics"
	mw "service/http/middleware"
	"service/logging"
	"service/metrics"
	serv "service/server"
)

const (
	// version     = "v1.0"
	serviceName = "template"
)

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
		logging.WithCaller(),
		logging.WithPID(),
	)

	logger.Info().Any("config", c).Send()

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
	forClose := RegisterMainServiceRoutes(mainRouter)

	defer func() {
		for _, closer := range forClose {
			err = closer.Close()
			if err != nil {
				logger.Err(err).Msg("failed to close:")
			}
		}
	}()

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

func RegisterMainServiceRoutes(r chi.Router) []io.Closer { //nolint:unparam
	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httpmetrics.RecordRequestHit)

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
