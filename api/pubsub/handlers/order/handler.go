package order

import (
	"github.com/rs/zerolog"
	"service/eserializer"
)

type Handler struct {
	logger       *zerolog.Logger
	deserializer eserializer.EventSerializer
}

func NewHandler(logger *zerolog.Logger, deserializer eserializer.EventSerializer) *Handler {
	return &Handler{logger: logger, deserializer: deserializer}
}
