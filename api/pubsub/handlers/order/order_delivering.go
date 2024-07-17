package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
	"service/domain/order/event"
)

func (h *Handler) OrderDelivering(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)
	eventOrderDelivering := &event.JSONDelivering{}

	err := h.unmarshaler.Unmarshal(msg.Payload, eventOrderDelivering)
	if err != nil {
		return errors.Wrap(err, "failed to parse order delivering event")
	}

	err = h.repository.Operate(ctx, eventOrderDelivering.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		delivering, err := stateOperator.DeliveringOrder()
		if err != nil || !delivering {
			return errors.Wrapf(err, "can`t set order`s state to delivering: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order delivering")
	}

	return nil
}
