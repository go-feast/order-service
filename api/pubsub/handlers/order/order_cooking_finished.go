package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
	"service/domain/order/event"
)

func (h *Handler) FinishedCooking(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)
	eventOrderFinished := &event.JSONOrderFinished{}

	err := h.unmarshaler.Unmarshal(msg.Payload, eventOrderFinished)
	if err != nil {
		return errors.Wrap(err, "failed to parse order finished cooking event")
	}

	err = h.repository.Operate(ctx, eventOrderFinished.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		finishedCooking, err := stateOperator.OrderFinished()
		if err != nil || !finishedCooking {
			return errors.Wrapf(err, "can`t set order`s state to finished cooking: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order finished cooking")
	}

	return nil
}
