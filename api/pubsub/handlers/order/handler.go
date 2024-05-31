package order

import (
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"service/eserializer"
)

type Handler struct {
	logger       *zerolog.Logger
	deserializer eserializer.EventSerializer
	tracer       trace.Tracer
}

func NewHandler(logger *zerolog.Logger, deserializer eserializer.EventSerializer, tracer trace.Tracer) *Handler {
	return &Handler{logger: logger, deserializer: deserializer, tracer: tracer}
}
