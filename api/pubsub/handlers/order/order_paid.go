package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
)

var _ message.HandlerFunc = (&Handler{}).OrderPaid

func (h *Handler) OrderPaid(msg *message.Message) ([]*message.Message, error) {
	eventOrderPaid := &order.JSONEventOrderPaid{}

	err := h.deserializer.Deserialize(msg.Payload, eventOrderPaid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse order paid event")
	}

	return nil, err
}
