package order

import "github.com/ThreeDotsLabs/watermill/message"

type Handler struct {
}

func (h *Handler) OrderCreated(_ *message.Message) ([]*message.Message, error) {
	return nil, nil
}
