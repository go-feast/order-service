package pubsub

import (
	"errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"gorm.io/gorm"
)

var (
	sqlPublisherConfig = sql.PublisherConfig{
		SchemaAdapter:        sql.DefaultPostgreSQLSchema{},
		AutoInitializeSchema: true,
	}
	sqlSubscriberConfig = sql.SubscriberConfig{
		SchemaAdapter:    sql.DefaultPostgreSQLSchema{},
		OffsetsAdapter:   sql.DefaultPostgreSQLOffsetsAdapter{},
		InitializeSchema: true,
	}
	kafkaTracer = kafka.NewOTELSaramaTracer()

	kafkaPublisherConfig = kafka.PublisherConfig{
		Marshaler: OTELMarshaler{},
		Tracer:    kafkaTracer,
	}

	kafkaSubscriberConfig = kafka.SubscriberConfig{
		Unmarshaler: kafka.DefaultMarshaler{},
		Tracer:      kafkaTracer,
	}
)

func NewSQLPublisher(db *gorm.DB, logger watermill.LoggerAdapter) (message.Publisher, error) {
	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}

	return sql.NewPublisher(sqldb, sqlPublisherConfig, logger)
}

func NewSQLSubscriber(db *gorm.DB, logger watermill.LoggerAdapter) (message.Subscriber, error) {
	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}

	return sql.NewSubscriber(sqldb, sqlSubscriberConfig, logger)
}

func NewKafkaPublisher(urls []string, logger watermill.LoggerAdapter) (message.Publisher, error) {
	if len(urls) == 0 {
		return nil, errors.New("must provide at least one publisher url")
	}

	kafkaPublisherConfig.Brokers = urls

	return kafka.NewPublisher(kafkaPublisherConfig, logger)
}

func NewKafkaSubscriber(urls []string, logger watermill.LoggerAdapter) (message.Subscriber, error) {
	if len(urls) == 0 {
		return nil, errors.New("must provide at least one publisher url")
	}

	kafkaSubscriberConfig.Brokers = urls

	return kafka.NewSubscriber(kafkaSubscriberConfig, logger)
}
