package config

import (
	"context"
	"fmt"
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

type ServiceConfig struct {
	DB           *DBConfig                `env:", prefix=POSTGRES_"`
	Redis        *RedisConfig             `env:", prefix=REDIS_"`
	Kafka        *KafkaConfig             `env:", prefix=KAFKA_"`
	Server       *MainServiceServerConfig `env:", prefix=SERVER_"`
	MetricServer *MetricServerConfig      `env:", prefix=METRICS_"`
	Environment  Environment              `env:"ENVIRONMENT,required"`
}

type ConsumerConfig struct {
	DB           *DBConfig           `env:", prefix=POSTGRES_"`
	Redis        *RedisConfig        `env:", prefix=REDIS_"`
	Kafka        *KafkaConfig        `env:", prefix=KAFKA_"`
	MetricServer *MetricServerConfig `env:", prefix=METRICS_"`
	Environment  Environment         `env:"ENVIRONMENT,required"`
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
	HOST     string `env:"HOST,required"`
	USER     string `env:"USER,required"`
	PASSWORD string `env:"PASSWORD,required"`
	DB       string `env:"DB,required"`
	SSL      string `env:"SSL,required"`
}

// DSN example: "postgres://username:password@localhost:5432/database_name?sslmode=disable"
func (c *DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		c.USER, c.PASSWORD, c.HOST, c.DB, c.SSL)
}

type KafkaConfig struct { //nolint:govet
	KafkaURL []string `env:"URL,required"`
}
type RedisConfig struct { //nolint:govet
	RedisURL string `env:"URL,required"`
}
