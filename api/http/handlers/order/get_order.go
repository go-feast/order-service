package order

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"net/http"
	"service/http/httpstatus"
)

type (
	GetOrderResponse struct {
	}
)

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	var (
		id  = chi.URLParam(r, "id")
		ctx = r.Context()
	)

	uuid, err := uuid.Parse(id)
	if err != nil {
		httpstatus.BadRequest(ctx, w, err)
		return
	}

	o, err := h.repository.Get(ctx, uuid)
	if err != nil {
		httpstatus.InternalServerError(ctx, w, err)
		return
	}

	response := o.ToEvent().ToJSON()

	render.JSON(w, r, response)
}
