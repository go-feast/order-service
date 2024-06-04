package logging

// TODO: ask (maybe change to zerolog)

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	"service/config"
	"time"
)

type OptionFunc func(zerolog.Context) zerolog.Context

// NewDefaultLogger returns a global log, which can be set once via New function.
// If New wasn't called at least once, zerolog.Nop will be returned
func NewDefaultLogger(w io.Writer) *zerolog.Logger {
	l := zerolog.New(w)
	return &l
}

func NewNopLogger() *zerolog.Logger {
	l := zerolog.Nop()
	return &l
}

var log *zerolog.Logger

// New initializes the log. Sets the global log once.
func New(opts ...OptionFunc) *zerolog.Logger {
	if log != nil {
		return log
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

	logger := ctx.Logger()

	for _, opt := range opts {
		logger.UpdateContext(opt)
	}

	log = &logger

	return log
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

type logEntry struct {
	l *zerolog.Logger
}

func (l *logEntry) NewLogEntry(_ *http.Request) middleware.LogEntry {
	return l
}

func (l *logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.l.Info().
		Int("status", status).
		Int("bytes", bytes).
		Any("header", header).
		Dur("elapsed", elapsed).
		Any("extra", extra).
		Send()
}

func (l *logEntry) Panic(v interface{}, _ []byte) {
	l.l.Panic().Any("value", v).Stack().Send()
}

func NewLogEntry() middleware.LogFormatter {
	return &logEntry{
		l: log,
	}
}
