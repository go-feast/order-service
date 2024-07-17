package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
	"service/domain/order/event"
)

var _ message.NoPublishHandlerFunc = ((*Handler)(nil)).OrderPaid

func (h *Handler) OrderPaid(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)

	eventOrderPaid := &event.JSONEventOrderPaid{}

	err := h.unmarshaler.Unmarshal(msg.Payload, eventOrderPaid)
	if err != nil {
		return errors.Wrap(err, "failed to parse order paid event")
	}

	err = h.repository.Operate(ctx, eventOrderPaid.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		orderPaid, payErr := stateOperator.PayOrder(eventOrderPaid.OrderID)
		if payErr != nil || !orderPaid {
			return errors.Wrapf(payErr, "can`t set order`s state to paid: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order paid")
	}

	return nil
}
