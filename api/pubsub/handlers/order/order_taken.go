package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
	"service/domain/order/event"
)

func (h *Handler) CookingTaken(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)
	eventOrderTaken := &event.JSONCourierTook{}

	err := h.unmarshaler.Unmarshal(msg.Payload, eventOrderTaken)
	if err != nil {
		return errors.Wrap(err, "failed to parse order taken event")
	}

	err = h.repository.Operate(ctx, eventOrderTaken.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		taken, err := stateOperator.CourierTookOrder(eventOrderTaken.CourierID)
		if err != nil || !taken {
			return errors.Wrapf(err, "can`t set order`s state to taken: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order taken")
	}

	return nil
}
