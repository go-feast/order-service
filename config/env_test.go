package config_test

import (
	"github.com/stretchr/testify/assert"
	"service/config"
	"testing"
)

const key = "ENVIRONMENT"

func TestEnvironment_String(t *testing.T) {
	const testString = "testString"
	e := config.Environment(testString)

	assert.NotSame(t, testString, e)
	assert.Equal(t, testString, e.String())
}

func TestMustGetEnvironment(t *testing.T) {
	value := "testing"

	t.Setenv(key, value)
	t.Run("assert successful", func(t *testing.T) {
		assert.NotPanics(t, func() {
			e := config.MustGetEnvironment()
			assert.Equal(t, e, config.Testing)
		})
	})
	t.Run("assert with different case", func(t *testing.T) {
		assert.NotPanics(t, func() {
			e := config.MustGetEnvironment()
			assert.Equal(t, e, config.Testing)
		})
	})
}

func TestSetGetEnv(t *testing.T) {
	t.Run("assert successful", func(t *testing.T) {
		v := config.Testing.String()
		fakeSetEnv(key, v)

		config.SetGetEnv(fakeGetEnv)

		assert.NotPanics(t, func() {
			e := config.MustGetEnvironment()
			assert.Equal(t, e, config.Testing)
		})
	})
}

func fakeGetEnv(key string) string { return testMap[key] }
func fakeSetEnv(key, v string)     { testMap[key] = v }

var testMap = map[string]string{}
