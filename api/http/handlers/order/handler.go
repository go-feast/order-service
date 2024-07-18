package order

import (
	"go.opentelemetry.io/otel/trace"
	"service/domain/order"
	"service/domain/shared/saver"
)

type Handler struct {
	tracer trace.Tracer

	saverService saver.Saver[*order.Order]

	// metrics

	// repositories eg.
	repository order.Repository
}

func NewHandler(
	tracer trace.Tracer,
	repository order.Repository,
	saverService saver.Saver[*order.Order],
) *Handler {
	return &Handler{
		tracer:       tracer,
		repository:   repository,
		saverService: saverService,
	}
}
