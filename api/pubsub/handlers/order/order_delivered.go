package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
	"service/domain/order/event"
)

func (h *Handler) OrderDelivered(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)
	eventOrderDelivered := &event.JSONDelivered{}

	err := h.unmarshaler.Unmarshal(msg.Payload, eventOrderDelivered)
	if err != nil {
		return errors.Wrap(err, "failed to parse order delivered event")
	}

	err = h.repository.Operate(ctx, eventOrderDelivered.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		delivered, err := stateOperator.OrderDelivered()
		if err != nil || !delivered {
			return errors.Wrapf(err, "can`t set order`s state to delivered: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order delivered")
	}

	return nil
}
