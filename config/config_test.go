package config_test

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"service/config"
	"testing"
	"time"
)

func TestParseEnvironment_with_prefix(t *testing.T) {
	t.Setenv("SERVER_DB_URL", "localhost:8080")
	t.Run("assert with prefix", func(t *testing.T) {
		c := &struct {
			config.DBConfig `env:",prefix=SERVER_DB_"`
		}{}

		err := config.ParseConfig(c)
		assert.NoError(t, err)
	})
	t.Run("assert without prefix", func(t *testing.T) {
		c := &config.DBConfig{}

		err := config.ParseConfig(c)
		assert.Error(t, err)
	})
}

func TestServerConfig(t *testing.T) {
	testCases := []struct {
		st any
	}{
		{&config.MainServiceServerConfig{}},
		{&config.MetricServerConfig{}},
	}

	for _, testCase := range testCases {
		tc := testCase
		v := reflect.ValueOf(tc.st)

		_, ok := v.Interface().(config.ServerConfig)
		assert.True(t, ok)
	}
}

func TestMainServiceServerConfig(t *testing.T) {
	c := config.MainServiceServerConfig{
		Port:         "80",
		Host:         "127.0.0.1",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	assert.Equal(t, c.HostPort(), "127.0.0.1:80")
	assert.Equal(t, c.WriteTimeoutDur(), 10*time.Second)
	assert.Equal(t, c.ReadTimeoutDur(), 10*time.Second)
	assert.Equal(t, c.IdleTimeoutDur(), 10*time.Second)
}

func TestMetricServerConfig(t *testing.T) {
	c := config.MetricServerConfig{
		Port:         "80",
		Host:         "127.0.0.1",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	assert.Equal(t, c.HostPort(), "127.0.0.1:80")
	assert.Equal(t, c.WriteTimeoutDur(), 10*time.Second)
	assert.Equal(t, c.ReadTimeoutDur(), 10*time.Second)
	assert.Equal(t, c.IdleTimeoutDur(), 10*time.Second)
}
