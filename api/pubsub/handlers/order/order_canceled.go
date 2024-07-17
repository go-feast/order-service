package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
	"service/domain/order/event"
)

func (h *Handler) OrderCanceled(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)

	eventOrderCanceled := &event.JSONCanceled{}

	err := h.unmarshaler.Unmarshal(msg.Payload, eventOrderCanceled)
	if err != nil {
		return errors.Wrap(err, "failed to parse order canceled event")
	}

	err = h.repository.Operate(ctx, eventOrderCanceled.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		canceled, cancelErr := stateOperator.CancelOrder(eventOrderCanceled.Reason)
		if cancelErr != nil || !canceled {
			return errors.Wrapf(cancelErr, "can`t set order`s state to canceled: order: %s", o.ID())
		}

		// process eventOrderCanceled.Reason

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order canceled")
	}

	return nil
}
