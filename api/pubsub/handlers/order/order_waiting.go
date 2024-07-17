package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
	"service/domain/order/event"
)

func (h *Handler) OrderWaitingForCourier(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)
	eventOrderWaiting := &event.JSONWaitingForCourier{}

	err := h.unmarshaler.Unmarshal(msg.Payload, eventOrderWaiting)
	if err != nil {
		return errors.Wrap(err, "failed to parse order waiting event")
	}

	err = h.repository.Operate(ctx, eventOrderWaiting.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		waiting, err := stateOperator.WaitForCourier()
		if err != nil || !waiting {
			return errors.Wrapf(err, "can`t set order`s state to waiting: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order waiting")
	}

	return nil
}
