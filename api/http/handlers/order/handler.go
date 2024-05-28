package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/render"
	"github.com/go-feast/topics"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"service/domain/order"
	"service/eserializer"
	"time"
)

type Handler struct {
	tracer trace.Tracer

	// publisher
	publisher message.Publisher

	serializer eserializer.EventSerializer

	// metrics

	// repositories eg.

}

func NewHandler() *Handler {
	return &Handler{}
}

type TakeOrderRequest struct {
	CustomerID   string   `json:"customer_id"`
	RestaurantID string   `json:"restaurant_id"`
	Meals        []string `json:"meals"`
	Destination  struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"destination"`
}

type TakeOrderResponse struct { //nolint:govet
	OrderID   string    `json:"order_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (h *Handler) TakeOrder(w http.ResponseWriter, r *http.Request) {
	var (
		_, span = h.tracer.Start(r.Context(), "take order")
	)

	defer span.End()

	takeOrder := &TakeOrderRequest{}

	err := render.DecodeJSON(r.Body, takeOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o, err := order.NewOrder(
		takeOrder.RestaurantID,
		takeOrder.CustomerID,
		takeOrder.Meals,
		takeOrder.Destination.Latitude,
		takeOrder.Destination.Longitude,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	JSONOrder := o.ToEvent().ToJSON()

	bytes, err := h.serializer.Serialize(JSONOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := message.NewMessage(uuid.NewString(), bytes)

	err = h.publisher.Publish(topics.OrderCreated.String(), msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := TakeOrderResponse{
		OrderID:   o.ID().String(),
		Timestamp: o.CreateAt(),
	}

	render.JSON(w, r, response)
}
