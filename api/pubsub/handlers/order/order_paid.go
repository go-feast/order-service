package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"service/domain/order"
)

var _ message.NoPublishHandlerFunc = ((*Handler)(nil)).OrderPaid

func (h *Handler) OrderPaid(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)

	eventOrderPaid := &order.JSONEventOrderPaid{}

	err := h.unmarshaler.Unmarshal(msg.Payload, eventOrderPaid)
	if err != nil {
		return errors.Wrap(err, "failed to parse order paid event")
	}

	err = h.repository.Operate(ctx, eventOrderPaid.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		orderPaid, err := stateOperator.PayOrder(eventOrderPaid.OrderID)
		if err != nil || !orderPaid {
			return errors.Wrapf(err, "can`t set order`s state to paid: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order paid")
	}

	return nil
}
