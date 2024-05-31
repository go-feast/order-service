package main

import (
	"context"
	"errors"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service/api/http/handlers/order"
	"service/closer"
	"service/config"
	"service/eserializer"
	mw "service/http/middleware"
	"service/logging"
	"service/metrics"
	"service/pubsub"
	serv "service/server"
	"service/tracing"
)

const (
	version     = "v1.0"
	serviceName = "order"
)

func main() {
	c := &config.ServiceConfig{}
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

	forClose := closer.NewCloser(logger)
	defer forClose.Close()

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	pubSubLogger := logging.NewPubSubLogger()

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   c.Kafka.KafkaURL,
			Marshaler: pubsub.OTELMarshaler{},
			Tracer:    kafka.NewOTELSaramaTracer(),
		},
		pubSubLogger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create kafka publisher")
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(version),
		),
		resource.WithProcess(),
		resource.WithOS(),
	)
	if err != nil {
		logger.Err(err).Msg("filed to create resource")
		return
	}

	if err = tracing.RegisterTracerProvider(ctx, res); err != nil {
		logger.Err(err).Msg("failed to register tracer provider")
		return
	}

	metrics.RegisterServiceName(serviceName)

	// main server
	mainServiceServer, mainRouter := serv.NewServer(c.Server)

	// metric server
	metricServer, metricRouter := serv.NewServer(c.MetricServer)

	// register routes
	//		main
	fc := RegisterMainServiceRoutes(mainRouter, publisher)

	forClose.AppendClosers(fc...)
	//		metric
	RegisterMetricRoute(metricRouter)

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
	r.Use(middleware.RequestLogger(logging.NewLogEntry()))
	r.Use(middleware.Recoverer)
}

func RegisterMainServiceRoutes(r chi.Router, pub message.Publisher) []io.Closer { //nolint:unparam
	// middlewares
	Middlewares(r)
	r.Get("/healthz", mw.Healthz)

	handler := order.NewHandler(
		otel.GetTracerProvider().Tracer(serviceName),
		pub,
		eserializer.JSONSerializer{},
	)

	r.With(mw.ResolveTraceIDInHTTP(serviceName)).
		Route("/api/v1", func(r chi.Router) {
			r.Route("/order", func(r chi.Router) {
				r.Post("/", handler.TakeOrder)
			})
		})

	return nil
}

func RegisterMetricRoute(r chi.Router) {
	handler := promhttp.Handler()
	r.Get("/metrics", handler.ServeHTTP)
}
