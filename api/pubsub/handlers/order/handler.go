package order

import (
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"service/domain/order"
	"service/event"
)

type Handler struct {
	logger      *zerolog.Logger
	unmarshaler event.Unmarshaler
	tracer      trace.Tracer
	repository  order.Repository
}

func NewHandler(
	logger *zerolog.Logger,
	unmarshaler event.Unmarshaler,
	tracer trace.Tracer,
	repository order.Repository,
) *Handler {
	return &Handler{
		logger:      logger,
		unmarshaler: unmarshaler,
		tracer:      tracer,
		repository:  repository,
	}
}
