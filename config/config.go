package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"net"
	"time"
)

func ParseConfig(v any) error {
	if err := envconfig.Process(context.TODO(), v); err != nil {
		return err
	}

	return nil
}

type ServerConfig interface {
	HostPort() string
	WriteTimeoutDur() time.Duration
	ReadTimeoutDur() time.Duration
	IdleTimeoutDur() time.Duration
	ReadHeaderTimeoutDur() time.Duration
}

type OTEL struct {
	TraceEndpoint string `env:"OTLP_TRACE_ENDPOINT"`
}

type Config struct {
	DB           *DBConfig                `env:", prefix=SERVER_DB_"`
	Redis        *RedisConfig             `env:", prefix=SERVER_REDIS_"`
	Rabbit       *RabbitMQConfig          `env:", prefix=SERVER_RABBITMQ_"`
	Server       *MainServiceServerConfig `env:", prefix=SERVER_"`
	MetricServer *MetricServerConfig      `env:", prefix=SERVER_METRICS_"`
	OTEL         *OTEL                    `env:", prefix=OTEL_EXPORTER_"`
	Environment  Environment              `env:"ENVIRONMENT,required"`
}

type MainServiceServerConfig struct { //nolint:govet
	Port         string        `env:"PORT,required"`
	Host         string        `env:"HOST,required"`
	WriteTimeout time.Duration `env:"WRITETIMEOUT,required"`
	ReadTimeout  time.Duration `env:"READTIMEOUT,required"`
	IdleTimeout  time.Duration `env:"IDLETIMEOUT,required"`
}

func (m *MainServiceServerConfig) HostPort() string {
	return net.JoinHostPort(m.Host, m.Port)
}

func (m *MainServiceServerConfig) WriteTimeoutDur() time.Duration {
	return m.WriteTimeout
}

func (m *MainServiceServerConfig) ReadTimeoutDur() time.Duration {
	return m.ReadTimeout
}

func (m *MainServiceServerConfig) IdleTimeoutDur() time.Duration {
	return m.IdleTimeout
}

func (m *MainServiceServerConfig) ReadHeaderTimeoutDur() time.Duration {
	return 0
}

type MetricServerConfig struct { //nolint:govet
	Port         string        `env:"PORT,required"`
	Host         string        `env:"HOST,required"`
	WriteTimeout time.Duration `env:"WRITETIMEOUT"`
	ReadTimeout  time.Duration `env:"READTIMEOUT"`
	IdleTimeout  time.Duration `env:"IDLETIMEOUT"`
}

func (m *MetricServerConfig) HostPort() string {
	return net.JoinHostPort(m.Host, m.Port)
}

func (m *MetricServerConfig) WriteTimeoutDur() time.Duration {
	return m.WriteTimeout
}

func (m *MetricServerConfig) ReadTimeoutDur() time.Duration {
	return m.ReadTimeout
}

func (m *MetricServerConfig) IdleTimeoutDur() time.Duration {
	return m.IdleTimeout
}

func (m *MetricServerConfig) ReadHeaderTimeoutDur() time.Duration {
	return 0
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
