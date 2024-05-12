package logging

// TODO: ask (maybe change to zerolog)

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"service/config"
	"time"
)

type OptionFunc func(zerolog.Context) zerolog.Context

// NewDefaultLogger returns a global logger, which can be set once via New function.
// If New wasn't called at least once, zerolog.Nop will be returned
func NewDefaultLogger(w io.Writer) *zerolog.Logger {
	l := zerolog.New(w)
	return &l
}

func NewNopLogger() *zerolog.Logger {
	l := zerolog.Nop()
	return &l
}

var logger *zerolog.Logger

// New initializes the logger. Sets the global logger once.
func New(opts ...OptionFunc) *zerolog.Logger {
	if logger != nil {
		return logger
	}

	var (
		out = os.Stdout
		ctx zerolog.Context
	)

	switch config.MustGetEnvironment() {
	case config.Production:
		ctx = zerolog.New(out).With().Timestamp().Caller()
	case config.Development, config.Local:
		ctx = zerolog.New(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}).With().Timestamp().Caller()
	case config.Testing:
		return NewNopLogger()
	}

	log := ctx.Logger()

	for _, opt := range opts {
		log.UpdateContext(opt)
	}

	return logger
}

func WithTimestamp() OptionFunc {
	return func(c zerolog.Context) zerolog.Context {
		zerolog.TimestampFieldName = "ts"
		return c.Timestamp()
	}
}

func WithServiceName(name string) OptionFunc {
	return func(c zerolog.Context) zerolog.Context {
		return c.Str("service.name", name)
	}
}

func WithPID() OptionFunc {
	return func(c zerolog.Context) zerolog.Context {
		return c.Int("pid", os.Getpid())
	}
}
