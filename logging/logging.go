package logging

// TODO: ask (maybe change to zerolog)

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"service/config"
	"time"
)

// Logger is a custom logger that logs messages using zerolog.Logger.
type Logger struct {
	logger *zerolog.Logger
}

func (l *Logger) UpdateContext(update func(c zerolog.Context) zerolog.Context) {
	l.logger.UpdateContext(update)
}
func (l *Logger) Trace() *zerolog.Event                        { return l.logger.Trace() }
func (l *Logger) Debug() *zerolog.Event                        { return l.logger.Debug() }
func (l *Logger) Info() *zerolog.Event                         { return l.logger.Info() }
func (l *Logger) Warn() *zerolog.Event                         { return l.logger.Warn() }
func (l *Logger) Error() *zerolog.Event                        { return l.logger.Error() }
func (l *Logger) Err(err error) *zerolog.Event                 { return l.logger.Err(err) }
func (l *Logger) Fatal() *zerolog.Event                        { return l.logger.Fatal() }
func (l *Logger) Panic() *zerolog.Event                        { return l.logger.Panic() }
func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event { return l.logger.WithLevel(level) }
func (l *Logger) Log() *zerolog.Event                          { return l.logger.Log() }
func (l *Logger) Print(v ...interface{})                       { l.logger.Print(v...) }
func (l *Logger) Printf(format string, v ...interface{})       { l.logger.Printf(format, v...) }
func (l *Logger) Println(v ...interface{})                     { l.logger.Println(v...) }

type OptionFunc func(zerolog.Context) zerolog.Context

// NewDefaultLogger returns a global logger, which can be set once via New function.
// If New wasn't called at least once, zerolog.Nop will be returned
func NewDefaultLogger(w io.Writer) *Logger {
	logger := zerolog.New(w)

	l := &Logger{logger: &logger}

	return l
}

func NewNopLogger() *Logger {
	logger := zerolog.Nop()

	l := &Logger{logger: &logger}

	return l
}

var logger *Logger

// New initializes the logger. Sets the global logger once.
func New(opts ...OptionFunc) *Logger {
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

	logger = &Logger{logger: &log}

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
