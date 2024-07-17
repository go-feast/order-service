package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
	"service/domain/order/event"
)

func (h *Handler) OrderClosed(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)
	eventOrderClosed := &event.JSONClosed{}

	err := h.unmarshaler.Unmarshal(msg.Payload, eventOrderClosed)
	if err != nil {
		return errors.Wrap(err, "failed to parse order closed event")
	}

	err = h.repository.Operate(ctx, eventOrderClosed.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		closed, err := stateOperator.CloseOrder()
		if err != nil || !closed {
			return errors.Wrapf(err, "can`t set order`s state to closed: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order closed")
	}

	return nil
}
