package config

import (
	"time"
)

type Environment string

func (e Environment) String() string { return string(e) }

const (
	// Production represents production environment
	Production Environment = "production"

	// Development represents developer environment
	Development Environment = "development"
)

type ServerConfig struct { //nolint:govet
	DBUrl        string
	AmmpqURL     string
	RedisURL     string
	Port         string
	Host         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
	Environment  Environment
}
