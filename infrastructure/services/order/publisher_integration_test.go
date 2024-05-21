//go:build integration

package order

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/topics"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"service/config"
	"testing"
	"time"
)

var cfg = struct {
	KafkaURL []string `env:"KAFKA_URL"`
}{}

func TestMain(m *testing.M) {
	// set-up kafka connection
	err := config.ParseConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// run tests
	os.Exit(m.Run())
}

func TestPublisherService_PublishOrderCreated_kafka(t *testing.T) {
	var (
		publisher, subscriber = testSetupKafkaPubSub(t)
		pubService            = NewPublisherService(publisher)
	)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	t.Run(
		"assert order created and unmarshaled in kafka",
		pubsubMessageSending(ctx, topics.OrderCreated.String(), pubService, subscriber),
	)
}

func testSetupKafkaPubSub(t *testing.T) (message.Publisher, message.Subscriber) {
	var (
		saramaSubscriberConfig = kafka.DefaultSaramaSubscriberConfig()
		logger                 = watermill.NopLogger{}
	)

	saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               cfg.KafkaURL,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: saramaSubscriberConfig,
			ConsumerGroup:         "test_consumer_group",
		},
		logger,
	)

	assert.NoError(t, err)

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   cfg.KafkaURL,
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)

	assert.NoError(t, err)

	return publisher, subscriber
}
