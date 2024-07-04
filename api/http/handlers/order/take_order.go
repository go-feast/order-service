package order

import (
	"github.com/go-chi/render"
	"net/http"
	"service/domain/order"
	"service/http/httpstatus"
	"time"
)

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
		ctx, span = h.tracer.Start(r.Context(), "take order")
	)

	defer span.End()

	takeOrder := &TakeOrderRequest{}

	err := render.DecodeJSON(r.Body, takeOrder)
	if err != nil {
		httpstatus.BadRequest(ctx, w, err)
		return
	}

	span.AddEvent("parsed TakeOrderRequest")

	o, err := order.NewOrder(
		takeOrder.RestaurantID,
		takeOrder.CustomerID,
		takeOrder.Meals,
		takeOrder.Destination.Latitude,
		takeOrder.Destination.Longitude,
	)
	if err != nil {
		httpstatus.BadRequest(ctx, w, err)
		return
	}

	err = h.saverService.Save(ctx, o)
	if err != nil {
		httpstatus.InternalServerError(ctx, w, err)
		return
	}

	span.AddEvent("created order")

	response := TakeOrderResponse{
		OrderID:   o.ID().String(),
		Timestamp: o.CreateAt(),
	}

	httpstatus.Created(w, response)
}
