package pubsub

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama" //nolint:staticcheck // module is deprecated
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type OTELMarshaler struct{}

type MetadataWrapper struct {
	message.Metadata
}

var _ kafka.Marshaler = OTELMarshaler{}

func (m MetadataWrapper) Keys() []string {
	keys := make([]string, 0, len(m.Metadata))
	for key := range m.Metadata {
		keys = append(keys, key)
	}

	return keys
}

func (o OTELMarshaler) Marshal(topic string, msg *message.Message) (*sarama.ProducerMessage, error) {
	if value := msg.Metadata.Get(kafka.UUIDHeaderKey); value != "" {
		return nil, errors.Errorf("metadata %s is reserved by watermill for message UUID", kafka.UUIDHeaderKey)
	}

	headers := []sarama.RecordHeader{{
		Key:   []byte(kafka.UUIDHeaderKey),
		Value: []byte(msg.UUID),
	}}
	for key, value := range msg.Metadata {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(value),
		})
	}

	producerMessage := &sarama.ProducerMessage{
		Topic:   topic,
		Value:   sarama.ByteEncoder(msg.Payload),
		Headers: headers,
	}

	otel.GetTextMapPropagator().
		Inject(msg.Context(), otelsarama.NewProducerMessageCarrier(producerMessage))

	return producerMessage, nil
}

func SpanFromMessage(msg *message.Message, traceName, spanName string, traceOpts []trace.TracerOption, spanOpts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx := msg.Context()
	extractedCtx := otel.GetTextMapPropagator().
		Extract(ctx, MetadataWrapper{msg.Metadata})

	return otel.GetTracerProvider().
		Tracer(traceName, traceOpts...).
		Start(extractedCtx, spanName, spanOpts...)
}
