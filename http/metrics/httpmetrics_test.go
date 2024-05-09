package httpmetrics_test

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"net/http"
	httpmetrics "service/http/metrics"
	"service/logging"
	"service/metrics"
	"testing"
)

func TestRegisterServer(t *testing.T) {
	t.Run("initial registry", func(t *testing.T) {
		assert.NotPanics(t, func() { httpmetrics.RegisterServer(nil, nil) })
	})
}

func TestUnregisterServer(t *testing.T) {
	t.Run("unregister after register", func(t *testing.T) {
		assert.NotPanics(t, func() { httpmetrics.UnregisterServer() })
	})
}

func TestHandler(t *testing.T) {
	t.Run("assert unregistered server should panic", func(t *testing.T) {
		httpmetrics.UnregisterServer()

		assert.Panics(t, func() { httpmetrics.Handler() })
	})

	t.Run("assert returning handler not nil", func(t *testing.T) {
		m := metrics.NewMetrics("test")
		mc := metrics.NewMetricCollector(m, prometheus.NewRegistry())

		logger := logging.New()

		httpmetrics.RegisterServer(mc, logger)

		var h http.HandlerFunc

		assert.NotPanics(t, func() { h = httpmetrics.Handler() })

		assert.HTTPStatusCode(t, h, http.MethodGet, "/", nil, http.StatusOK)
	})
}

func TestRecordRequestHit(t *testing.T) {
	m := metrics.NewMetrics("test")
	mc := metrics.NewMetricCollector(m, prometheus.NewRegistry())

	logger := logging.New()

	httpmetrics.RegisterServer(mc, logger)

	var h http.HandlerFunc = func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }

	next := httpmetrics.RecordRequestHit(h)

	assert.HTTPStatusCode(t, next.ServeHTTP, http.MethodGet, "/", nil, http.StatusOK)
}
