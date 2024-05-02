package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
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
	version     = "v1.0"
	serviceName = "template"
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

	// metric server
	metricServer, metricRouter := serv.NewServer(c.MetricServer)
	httpmetrics.RegisterServer(
		metricServer,
		metrics.NewMetricsService(metrics.NewMetrics(serviceName+version), prometheus.NewRegistry()),
		logger,
	)

	// main server
	mainServiceServer, mainRouter := serv.NewServer(c.Server)

	// register routes
	//		main
	forClose := RegisterMainServiceRoutes(mainRouter)

	defer func() {
		for _, closer := range forClose {
			err := closer.Close()
			if err != nil {
				logger.Error("failed to close:", zap.Error(err))
			}
		}
	}()

	// 		metrics
	RegisterMetricRoutes(metricRouter)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()
	if err = serv.Run(ctx, logger, mainServiceServer, metricServer); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("app", zap.Error(err))
	}
}

func RegisterMainServiceRoutes(r chi.Router) []io.Closer {
	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httpmetrics.RecordRequestHit)

	// routes
	//forClose := make([]io.Closer, )

	return nil
}

// RegisterMetricRoutes initialize routes for metrics. Example:
//
// [/metric] - for accessing metrics
//
// [/ping] [/healthz] [/readyz] - for checking if service alive
func RegisterMetricRoutes(r chi.Router) {
	r.Get("/healthz", mw.Healthz)
	r.Get("/readyz", mw.Readyz)
	r.Get("/ping", mw.Ping)

	r.Get("/metrics", promhttp.Handler().ServeHTTP)
}
