package order

import (
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"service/event"
)

type Handler struct {
	logger      *zerolog.Logger
	unmarshaler event.Unmarshaler
	tracer      trace.Tracer
}

func NewHandler(logger *zerolog.Logger, unmarshaler event.Unmarshaler, tracer trace.Tracer) *Handler {
	return &Handler{logger: logger, unmarshaler: unmarshaler, tracer: tracer}
}
