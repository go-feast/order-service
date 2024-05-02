package metrics_test

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"service/metrics"
	"testing"
)

const _test_string = "test"

func TestNewMetrics(t *testing.T) {
	m := metrics.NewMetrics(_test_string)

	assert.NotNil(t, m.RequestsHit)
	assert.NotNil(t, m.RequestProceedingDuration)

	m2 := metrics.NewMetrics(_test_string)

	assert.NotSame(t, m, m2)
}

func TestMetrics_Collectors(t *testing.T) {
	testString := _test_string
	m := metrics.NewMetrics(testString)

	v := reflect.ValueOf(metrics.Metrics{})

	assert.Equal(t, len(m.Collectors()), v.NumField())
}
