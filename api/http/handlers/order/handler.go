package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/trace"
	"service/event"
)

type Handler struct {
	tracer trace.Tracer

	publisher message.Publisher

	marshaler event.Marshaler

	// metrics

	// repositories eg.
}

func NewHandler(
	tracer trace.Tracer,
	publisher message.Publisher,
	marshaler event.Marshaler,
) *Handler {
	return &Handler{
		tracer:    tracer,
		publisher: publisher,
		marshaler: marshaler,
	}
}
