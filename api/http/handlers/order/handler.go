package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"service/eserializer"
)

type Handler struct {
	_ trace.Tracer

	// publisher
	_ message.Publisher

	_ eserializer.SerializeDeserializer

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
