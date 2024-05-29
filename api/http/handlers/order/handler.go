package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/trace"
	"service/eserializer"
)

type Handler struct {
	tracer trace.Tracer

	publisher message.Publisher

	serializer eserializer.EventSerializer

	// metrics

	// repositories eg.
}

func NewHandler(
	tracer trace.Tracer,
	publisher message.Publisher,
	serializer eserializer.EventSerializer,
) *Handler {
	return &Handler{
		tracer:     tracer,
		publisher:  publisher,
		serializer: serializer,
	}
}
