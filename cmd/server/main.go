package main

import (
	"context"
	"errors"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"log"
	"net/http"
	"service/config"
	"service/logging"
)

func main() {
	// config
	c := &config.ServerConfig{}

	err := envconfig.Process("server", c)
	if err != nil {
		log.Fatal(err)
	}

	// logger
	logger, err := logging.NewLogger(c.Environment,
		logging.WithOutputPaths("stdout"),
		logging.WithErrorOutputPaths("stderr"),
	)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("config", zap.Any("config", c))

	// app
	ctx := context.Background()
	if err = app(ctx, c, logger); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("app", zap.Error(err))
	}
}
