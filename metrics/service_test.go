package metrics_test

import (
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"service/metrics"
	"testing"
)

func TestNewMetricsService(t *testing.T) {
	m := metrics.NewMetrics(_test_string)

	ms1 := metrics.NewMetricsService(m, nil)
	ms2 := metrics.NewMetricsService(m, nil)

	assert.NotSame(t, ms1, ms2)
}

func TestMetricService_Handler(t *testing.T) {
	m := metrics.NewMetrics("example_2")

	ms := metrics.NewMetricsService(m, prometheus.NewRegistry())

	assert.NotNil(t, ms.Handler())
}

func assertRequestsHitEqual(t *testing.T, metric *metrics.Metrics, expected float64) {
	m := &io_prometheus_client.Metric{}
	metric.RequestsHit.WithLabelValues("GET", "/").Write(m)
	assert.Equal(t, expected, *m.Counter.Value)
}

func TestMetricService_RecordRequestHit(t *testing.T) {
	metric := metrics.NewMetrics(_test_string + "2")

	ms := metrics.NewMetricsService(metric, prometheus.NewRegistry())
	assertRequestsHitEqual(t, metric, 0.0)

	ms.RecordRequestHit("GET", "/")

	assertRequestsHitEqual(t, metric, 1.0)
}
