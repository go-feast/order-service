package metrics_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"service/metrics"
	"testing"
)

func TestNewMetrics(t *testing.T) {
	const serviceName = "testing"

	t.Run("assert metrics.NewMetrics returns not nil metrics", func(t *testing.T) {
		m := metrics.NewMetrics(serviceName)
		require.NotNil(t, m)

		assertFieldsNotNil(*m)
	})
	t.Run("assert assertFieldsnotNil panics when field is nil", func(t *testing.T) {
		m := metrics.NewMetrics(serviceName)

		testStruct := struct {
			*metrics.Metrics
			v *int
		}{m, nil} // *int should be nil, to test if function panics

		assert.Panics(t, func() {
			assertFieldsNotNil(testStruct)
		})
	})
}

func assertFieldsNotNil[T comparable](v T) {
	refV := reflect.ValueOf(v)

	for i := 0; i < refV.NumField(); i++ {
		field := refV.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			panic("field is nil")
		}
	}
}
