package main

import (
	"context"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-feast/topics"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service/api/pubsub/handlers/order"
	"service/closer"
	"service/config"
	"service/event"
	mw "service/http/middleware"
	repository "service/infrastructure/repositories/order/gorm"
	"service/logging"
	"service/metrics"
	"service/pubsub"
	serv "service/server"
	"service/tracing"
)

const (
	version     = "v1.0"
	serviceName = "order_consumer"
)

func main() {
	c := &config.ConsumerConfig{}

	err := config.ParseConfig(c)
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.New(
		logging.WithServiceName(serviceName),
		logging.WithTimestamp(),
		logging.WithPID(),
	)

	logger.Info().Any("config", c).Send()

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

	Closer := closer.NewCloser(logger)
	defer Closer.Close()

	// metric server
	metricServer, metricRouter := serv.NewServer(c.MetricServer)

	RegisterMetricRoute(metricRouter)

	pubSubLogger := logging.NewWatermillAdapter()
	// consumer router
	router, err := message.NewRouter(message.RouterConfig{}, pubSubLogger)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to create message router")
	}

	Closer.AppendClosers(closer.C{Name: "router", Closer: router})

	driverName := "pgx/v5"

	db, err := gorm.Open(postgres.New(
		postgres.Config{
			DriverName: driverName,
			DSN:        c.DB.DSN(),
		}), &gorm.Config{})
	if err != nil {
		logger.Fatal().Err(err).
			Str("dsn", c.DB.DSN()).
			Str("driver", driverName).
			Msg("failed to connect to database")
	}

	closers := RegisterConsumerHandlers(router, db, c.Kafka)

	Closer.AppendClosers(closers...)

	go func() {
		e := router.Run(ctx)
		if e != nil {
			logger.Error().Err(err).Msgf("failed to run consumer: %s", e.Error())
			return
		}

		logger.Info().Msg("exiting consumer")
	}()

	_, errCh := serv.Run(ctx, metricServer)

	for err = range errCh {
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Err(err).Send()
		}
	}
}

func RegisterMetricRoute(r chi.Router) {
	handler := promhttp.Handler()
	r.Get("/metrics", handler.ServeHTTP)
	r.Get("/healthz", mw.Healthz)
}

func RegisterConsumerHandlers(r *message.Router, db *gorm.DB, c *config.KafkaConfig) []closer.C {
	subscriberSQL, err := pubsub.NewSQLSubscriber(db, logging.NewWatermillAdapter())
	if err != nil {
		panic(err)
	}

	publisherKafka, err := pubsub.NewKafkaPublisher(c.KafkaURL, logging.NewWatermillAdapter())
	if err != nil {
		panic(err)
	}

	subscriberKafka, err := pubsub.NewKafkaSubscriber(c.KafkaURL, logging.NewWatermillAdapter())
	if err != nil {
		panic(err)
	}

	orderRepository := repository.NewOrderRepository(db)

	handler := order.NewHandler(
		logging.New(),
		event.JSONMarshaler{},
		otel.GetTracerProvider().Tracer(serviceName),
		orderRepository,
	)

	r.AddHandler(
		"handler.order.created",
		topics.OrderCreated.String(),
		subscriberSQL,
		topics.OrderCreated.String(),
		publisherKafka,
		handler.OrderCreated,
	)

	registerOrderStateHandlers(r, handler, subscriberKafka)
	return []closer.C{
		{Name: "kafka pub", Closer: publisherKafka},
		{Name: "sql sub", Closer: subscriberSQL},
		{Name: "kafka sub", Closer: subscriberKafka},
	}
}

func registerOrderStateHandlers(r *message.Router, handler *order.Handler, subKafka message.Subscriber) {

	r.AddNoPublisherHandler(
		"handler.order.paid",
		topics.Paid.String(),
		subKafka,
		handler.OrderPaid,
	)

	r.AddNoPublisherHandler(
		"order.cooking",
		topics.Cooking.String(),
		subKafka,
		handler.CookingOrder,
	)

	r.AddNoPublisherHandler(
		"order.cooking.finished",
		topics.CookingFinished.String(),
		subKafka,
		handler.FinishedCooking,
	)

	r.AddNoPublisherHandler(
		"order.waiting",
		topics.WaitingForCourier.String(),
		subKafka,
		handler.OrderWaitingForCourier,
	)

	r.AddNoPublisherHandler(
		"order.taken",
		topics.CourierTook.String(),
		subKafka,
		handler.CookingTaken,
	)

	r.AddNoPublisherHandler(
		"order.delivering",
		topics.Delivering.String(),
		subKafka,
		handler.OrderDelivering,
	)

	r.AddNoPublisherHandler(
		"order.delivered",
		topics.Delivered.String(),
		subKafka,
		handler.OrderDelivered,
	)

	r.AddNoPublisherHandler(
		"order.closed",
		topics.Closed.String(),
		subKafka,
		handler.OrderClosed,
	)

	r.AddNoPublisherHandler(
		"order.canceled",
		topics.Canceled.String(),
		subKafka,
		handler.OrderCanceled,
	)

}
