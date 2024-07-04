package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"service/event"
)

type Handler struct {
	logger      *zerolog.Logger
	unmarshaler event.Unmarshaler
	tracer      trace.Tracer
	publisher   message.Publisher
}

func NewHandler(
	logger *zerolog.Logger,
	unmarshaler event.Unmarshaler,
	tracer trace.Tracer,
	publisher message.Publisher,
) *Handler {
	return &Handler{
		logger:      logger,
		unmarshaler: unmarshaler,
		tracer:      tracer,
		publisher:   publisher,
	}
}
