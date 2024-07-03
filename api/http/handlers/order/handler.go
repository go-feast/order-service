package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/trace"
	"service/domain/order"
	"service/event"
)

type Handler struct {
	tracer trace.Tracer

	publisher message.Publisher

	marshaler event.Marshaler

	// metrics

	// repositories eg.
	repository order.Repository
}

func NewHandler(
	tracer trace.Tracer,
	publisher message.Publisher,
	marshaler event.Marshaler,
	repository order.Repository,
) *Handler {
	return &Handler{
		tracer:     tracer,
		publisher:  publisher,
		marshaler:  marshaler,
		repository: repository,
	}
}
