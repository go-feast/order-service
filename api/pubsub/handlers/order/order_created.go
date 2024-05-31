package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"service/pubsub"
)

func (h *Handler) OrderCreated(msg *message.Message) error {
	_, span := pubsub.SpanFromMessage(
		msg,
		"consumer.order.created",
		"order.created handler",
		nil,
	)
	defer span.End()

	h.logger.Info().Str("msg-id", msg.UUID).Msg("Received message from topic OrderCreated")

	return nil
}
