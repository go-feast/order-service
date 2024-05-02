package metrics_test

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"service/metrics"
	"testing"
)

const testString = "test"

func TestNewMetrics(t *testing.T) {
	m := metrics.NewMetrics(testString)

	assert.NotNil(t, m.RequestsHit)
	assert.NotNil(t, m.RequestProceedingDuration)

	m2 := metrics.NewMetrics(testString)

	assert.NotSame(t, m, m2)
}

func TestMetrics_Collectors(t *testing.T) {
	testString := testString
	m := metrics.NewMetrics(testString)

	v := reflect.ValueOf(metrics.Metrics{})

	assert.Equal(t, len(m.Collectors()), v.NumField())
}
