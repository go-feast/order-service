package httpmetrics_test

import (
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"service/config"
	httpmetrics "service/http/metrics"
	"service/logging"
	"service/metrics"
	"testing"
)

func TestRegisterServer(t *testing.T) {
	metric := metrics.NewMetrics("example1")
	createServer(metric)

	handler := httpmetrics.RecordRequestHit(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	assertCounterMetricEqual(t, metric, 0.0)

	assert.HTTPStatusCode(t, handler.ServeHTTP, http.MethodGet, "/", nil, http.StatusOK)

	assertCounterMetricEqual(t, metric, 1.0)

	createServer(metric)

	assertCounterMetricEqual(t, metric, 1.0)
}

func assertCounterMetricEqual(t *testing.T, metric *metrics.Metrics, expected float64) {
	m := &io_prometheus_client.Metric{}
	metric.RequestsHit.WithLabelValues("GET", "/").Write(m) //nolint:errcheck
	assert.Equal(t, expected, *m.Counter.Value)
}

func TestRecordRequestHit(t *testing.T) {
	metric := metrics.NewMetrics("example2")
	createServer(metric)

	handler := httpmetrics.RecordRequestHit(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	assertCounterMetricEqual(t, metric, 0.0)

	assert.HTTPStatusCode(t, handler.ServeHTTP, http.MethodGet, "/", nil, http.StatusOK)
	// 1.0 throws an error while processing with other test, but when executed solo, 1.0 passes
	assertCounterMetricEqual(t, metric, 1.0)
}

func createServer(metric *metrics.Metrics) {
	server := &http.Server{
		Addr: ":8080",
	}

	l, _ := logging.NewLogger(config.Development)
	service := metrics.NewMetricsService(metric, prometheus.NewRegistry())
	httpmetrics.RegisterServer(server, service, l)
}
