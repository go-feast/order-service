package order

import "github.com/ThreeDotsLabs/watermill/message"

func (h *Handler) OrderCreated(msg *message.Message) error {
	h.logger.Info().Str("msg-id", msg.UUID).Msg("Received message from topic OrderCreated")

	return nil
}
