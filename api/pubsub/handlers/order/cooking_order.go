package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
	"service/domain/order/event"
)

func (h *Handler) CookingOrder(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)
	eventOrderCooking := &event.JSONOrderCooking{}

	err := h.unmarshaler.Unmarshal(msg.Payload, eventOrderCooking)
	if err != nil {
		return errors.Wrap(err, "failed to parse order cooking event")
	}

	err = h.repository.Operate(ctx, eventOrderCooking.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		cooking, err := stateOperator.CookOrder()
		if err != nil || !cooking {
			return errors.Wrapf(err, "can`t set order`s state to cooking: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order cooking")
	}

	return nil
}
