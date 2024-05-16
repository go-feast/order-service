package logging

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog"
)

func NewPubSubLogger() *PubSubLogger {
	return &PubSubLogger{l: logger}
}

type PubSubLogger struct {
	l *zerolog.Logger
}

func (p *PubSubLogger) Error(msg string, err error, fields watermill.LogFields) {
	p.l.Error().Err(err).Any("fields", fields).Msg(msg)
}

func (p *PubSubLogger) Info(msg string, fields watermill.LogFields) {
	p.l.Info().Any("fields", fields).Msg(msg)
}

func (p *PubSubLogger) Debug(msg string, fields watermill.LogFields) {
	p.l.Debug().Any("fields", fields).Msg(msg)
}

func (p *PubSubLogger) Trace(msg string, fields watermill.LogFields) {
	p.l.Trace().Any("fields", fields).Msg(msg)
}

func (p *PubSubLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	p.l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Any("fields", fields)
	})

	return p
}
