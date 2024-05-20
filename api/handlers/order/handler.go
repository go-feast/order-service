package order

import (
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"service/domain/order"
)

type Handler struct {
	_ trace.Tracer

	// publisher
	_ order.Publisher

	// metrics

	// repositories eg.

}

func NewHandler() *Handler {
	return &Handler{}
}

type TakeOrderRequest struct {
}

type TakeOrderResponse struct {
}

func (h *Handler) TakeOrder(_ http.ResponseWriter, _ *http.Request) {

}
