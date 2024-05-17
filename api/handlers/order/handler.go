package order

import (
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type Handler struct {
	_ trace.Tracer
	// metrics

	// repositories eg.
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) TakeOrder(_ http.ResponseWriter, _ *http.Request) {

}
