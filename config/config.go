package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
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

type Config struct {
	DB           *DBConfig           `env:", prefix=SERVER_DB_"`
	Redis        *RedisConfig        `env:", prefix=SERVER_REDIS_"`
	Rabbit       *RabbitMQConfig     `env:", prefix=SERVER_RABBITMQ_"`
	Server       *ServerConfig       `env:", prefix=SERVER_"`
	MetricServer *MetricServerConfig `env:", prefix=SERVER_METRICS_"`
}

type ServerConfig struct { //nolint:govet
	Port         string        `env:"PORT,required"`
	Host         string        `env:"HOST,required"`
	WriteTimeout time.Duration `env:"WRITETIMEOUT,required"`
	ReadTimeout  time.Duration `env:"READTIMEOUT,required"`
	IdleTimeout  time.Duration `env:"IDLETIMEOUT,required"`
	Environment  Environment   `env:"ENVIRONMENT,required"`
}

type MetricServerConfig struct { //nolint:govet
	Port         string        `env:"PORT,required"`
	Host         string        `env:"HOST,required"`
	WriteTimeout time.Duration `env:"WRITETIMEOUT"`
	ReadTimeout  time.Duration `env:"READTIMEOUT"`
	IdleTimeout  time.Duration `env:"IDLETIMEOUT"`
}

type DBConfig struct { //nolint:govet
	DBURL string `env:"URL,required"`
}
type RabbitMQConfig struct { //nolint:govet
	RabbitMQURL string `env:"URL,required"`
}
type RedisConfig struct { //nolint:govet
	RedisURL string `env:"URL,required"`
}

func ParseEnvironment(v any) error {
	if err := envconfig.Process(context.TODO(), v); err != nil {
		return err
	}

	return nil
}
