package main

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service/closer"
	"service/config"
	mw "service/http/middleware"
	"service/logging"
	"service/metrics"
	serv "service/server"
	"service/tracing"
)

const (
	serviceName = "template-consumer"
)

func main() {
	c := &config.ConsumerConfig{}

	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

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

	if err = tracing.RegisterTracerProvider(ctx, serviceName); err != nil {
		logger.Err(err).Msg("failed to register tracer provider")
		return
	}

	metrics.RegisterServiceName(serviceName)

	Closer := closer.NewCloser(logger)
	defer Closer.Close()

	// metric server
	metricServer, metricRouter := serv.NewServer(c.MetricServer)

	RegisterMetricRoute(metricRouter)

	pubSubLogger := logging.NewPubSubLogger()
	// consumer router
	router, err := message.NewRouter(message.RouterConfig{}, pubSubLogger)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to create message router")
	}

	Closer.AppendClosers(router)

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               c.Kafka.KafkaURL,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         "test_consumer_group",
		},
		pubSubLogger,
	)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to create kafka subscriber")
	}

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   c.Kafka.KafkaURL,
			Marshaler: kafka.DefaultMarshaler{},
		},
		pubSubLogger,
	)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to create kafka publisher")
	}

	closers := RegisterConsumerHandlers(router, subscriber, publisher)

	Closer.AppendClosers(closers...)

	go func() {
		e := router.Run(ctx)
		if e != nil {
			logger.Error().Err(err).Msg("failed to run consumer")
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

func RegisterConsumerHandlers(_ *message.Router, subscriber message.Subscriber, publisher message.Publisher) []io.Closer {
	return []io.Closer{
		subscriber,
		publisher,
	}
}
