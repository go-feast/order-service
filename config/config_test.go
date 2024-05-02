package config_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"service/config"
	"strings"
	"testing"
)

func TestParseEnvironment_with_prefix(t *testing.T) {
	t.Setenv("SERVER_DB_URL", "localhost:8080")

	c := &struct {
		config.DBConfig `env:",prefix=SERVER_DB_"`
	}{}

	err := config.ParseEnvironment(c)
	assert.NoError(t, err)
}

func TestParseEnvironment_without_prefix(t *testing.T) {
	t.Setenv("SERVER_DB_URL", "localhost:8080")

	c := &config.DBConfig{}

	err := config.ParseEnvironment(c)
	assert.NotNil(t, err)
}

func TestEnvironment_String(t *testing.T) {
	var e config.Environment = "f"

	assert.Equal(t, "f", e.String())
}

func TestEnvironment_Production_and_Develop(t *testing.T) {
	require.Equal(t, strings.ToLower(config.Production.String()), "production")
	require.Equal(t, strings.ToLower(config.Development.String()), "development")
}

func TestMainServiceServerConfig(t *testing.T) {
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
