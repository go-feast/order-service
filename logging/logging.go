package logging

// TODO: ask (maybe change to zerolog)

import (
	"go.uber.org/zap"
	"service/config"
)

// OptionFunc represents a function that can be used to configure the logger.
type OptionFunc func(*zap.Config)

// NewLogger initializes the logger based on the provided environment and options.
func NewLogger(options ...OptionFunc) (*zap.Logger, error) {
	var (
		err error
		c   zap.Config
	)

	switch config.MustGetEnvironment() {
	case config.Production:
		c = zap.NewProductionConfig()
		// FIXME: refactor
	case config.Local:
		c = zap.NewDevelopmentConfig()
		// FIXME: refactor
	case config.Testing:
		c = zap.NewDevelopmentConfig()
	case config.Development:
		c = zap.NewDevelopmentConfig()

	default:
	}

	for _, opt := range options {
		opt(&c)
	}

	logger, err := c.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// WithLevel sets the log level of the logger.
func WithLevel(level zap.AtomicLevel) OptionFunc {
	return func(config *zap.Config) {
		config.Level = level
	}
}

// WithOutputPaths sets the output paths of the logger.
func WithOutputPaths(paths ...string) OptionFunc {
	return func(config *zap.Config) {
		config.OutputPaths = paths
	}
}

// WithErrorOutputPaths sets the error output paths of the logger.
func WithErrorOutputPaths(paths ...string) OptionFunc {
	return func(config *zap.Config) {
		config.ErrorOutputPaths = paths
	}
}
