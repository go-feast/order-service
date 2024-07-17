package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service/api/http/handlers/order"
	"service/closer"
	"service/config"
	domain "service/domain/order"
	"service/event"
	mw "service/http/middleware"
	"service/infrastructure/outbox"
	repository "service/infrastructure/repositories/order/gorm"
	"service/logging"
	"service/metrics"
	"service/pubsub"
	serv "service/server"
	"service/tracing"
)

const (
	version     = "v1.0"
	serviceName = "order"
	driverName  = "pgx/v5"
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

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName:           "pgx",
		DSN:                  c.DB.DSN(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
		return
	}

	db = db.WithContext(ctx)

	domain.InitializeOrderScheme(db)
  
	// main server
	mainServiceServer, mainRouter := serv.NewServer(c.Server)

	// metric server
	metricServer, metricRouter := serv.NewServer(c.MetricServer)

	// register routes
	//		main
	fc := RegisterMainServiceRoutes(mainRouter, db)

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

func RegisterMainServiceRoutes(
	r chi.Router,
	db *gorm.DB,
) []closer.C { //nolint:unparam

	// middlewares
	Middlewares(r)
	r.Get("/healthz", mw.Healthz)

	publisher, err := pubsub.NewSQLPublisher(db, logging.NewWatermillAdapter())
	if err != nil {
		panic("Failed to create publisher")
	}

	orderRepository := repository.NewOrderRepository(db)
	saver := outbox.NewOutbox(
		publisher,
		orderRepository,
		event.JSONMarshaler{},
	)

	handler := order.NewHandler(
		otel.GetTracerProvider().Tracer(serviceName),
		orderRepository,
		saver,
	)

	r.With(mw.ResolveTraceIDInHTTP(serviceName)).
		Route("/api/v1", func(r chi.Router) {
			r.Route("/order", func(r chi.Router) {
				r.Post("/", handler.TakeOrder)
			})
		})

	return []closer.C{
		{Name: "pub", Closer: publisher},
	}
}

func RegisterMetricRoute(r chi.Router) {
	handler := promhttp.Handler()
	r.Get("/metrics", handler.ServeHTTP)
}
