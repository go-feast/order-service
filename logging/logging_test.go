package logging_test

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"reflect"
	"service/logging"
	"testing"
)

func TestOutputPaths(t *testing.T) {
	t.Parallel()
	t.Run("output paths", TestWithOutputPaths)
	t.Run("output error paths", TestWithErrorOutputPaths)
	t.Run("level", TestWithLevel)
	t.Run("new logger", TestNewLoggerDevNProdAreNotEqual)
}

func TestNewLoggerDevNProdAreNotEqual(t *testing.T) {
	dev, _ := logging.NewLogger()
	prod, _ := logging.NewLogger()

	assert.Equal(t, false, reflect.DeepEqual(dev, prod))
}

func TestWithLevel(t *testing.T) {
	var testLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

	c := &zap.Config{}
	paths := logging.WithLevel(testLevel)
	paths(c)
	assert.Equal(t, testLevel, c.Level)
}

func TestWithOutputPaths(t *testing.T) {
	var testPaths = []string{"test", "test2", "test3"}

	c := &zap.Config{}
	paths := logging.WithOutputPaths(testPaths...)
	paths(c)
	assert.Equal(t, testPaths, c.OutputPaths)
}

func TestWithErrorOutputPaths(t *testing.T) {
	var testPaths = []string{"test", "test2", "test3"}

	c := &zap.Config{}
	paths := logging.WithErrorOutputPaths(testPaths...)
	paths(c)
	assert.Equal(t, testPaths, c.ErrorOutputPaths)
}
