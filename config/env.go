package config

import (
	"errors"
	"os"
)

type Environment string

func (e Environment) String() string { return string(e) }

const (
	// Production represents production environment
	Production Environment = "production"

	// Development represents developer environment
	Development Environment = "development"

	// Local represents developer environment
	Local Environment = "local"

	// Testing represents developer environment
	Testing Environment = "testing"
)

var envs = map[string]Environment{
	"production":  Production,
	"development": Development,
	"local":       Local,
	"testing":     Testing,
}

var getenv = os.Getenv

func SetGetEnv(f func(string) string) {
	getenv = f
}

var ErrPanicMsg = errors.New("invalid environment")

// MustGetEnvironment returns environment variable by specified key via getenv function.
func MustGetEnvironment() Environment {
	const key = "ENVIRONMENT"
	str := getenv(key)

	if env, ok := envs[str]; ok {
		return env
	}

	panic(ErrPanicMsg)
}
