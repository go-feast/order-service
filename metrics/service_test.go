package metrics_test

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"service/metrics"
	"sync"
	"testing"
	"time"
)

func TestMetrics_Collectors(t *testing.T) {
	t.Run("test metrics.Meteics is nil", func(t *testing.T) {
		mc := metrics.NewMetricCollector(nil, nil)

		assert.Nil(t, mc.Collectors())
		assert.Len(t, mc.Collectors(), 0)
	})
}

func TestMetricCollector(t *testing.T) {
	t.Run("assert metrics.NewMetricCollector returns not nil", func(t *testing.T) {
		collector := metrics.NewMetricCollector(nil, nil)
		assert.NotNil(t, collector)
	})
}

func TestMetricCollector_Handler(t *testing.T) {
	const testString = "testString"

	t.Parallel()

	t.Run("assert (metrics=initialized, reg=initialized)", func(t *testing.T) {
		m := metrics.NewMetrics(testString)
		collector := metrics.NewMetricCollector(m, prometheus.NewRegistry())

		assert.NotPanics(t, func() { collector.Handler() }) //nolint:wsl
	})

	t.Run("assert (metrics=nil, reg=initialized)", func(t *testing.T) {
		collector := metrics.NewMetricCollector(nil, prometheus.NewRegistry())

		var h http.HandlerFunc

		assert.NotPanics(t, func() { h = collector.Handler() })
		assert.NotNil(t, h)
		assert.HTTPStatusCode(t, h, http.MethodGet, "/", nil, http.StatusOK)
	})

	t.Run("assert (metrics=nil, reg=nil)", func(t *testing.T) {
		collector := metrics.NewMetricCollector(nil, nil)

		assert.Panics(t, func() { collector.Handler() })
	})
	// TODO: ask
	t.Run("assert handler contains our metrics", func(t *testing.T) {
		m := metrics.NewMetrics(testString)
		collector := metrics.NewMetricCollector(m, prometheus.NewRegistry())

		h := collector.Handler()

		const metricHTTPResponse = `# HELP promhttp_metric_handler_errors_total Total number of internal errors encountered by the promhttp metric handler.
# TYPE promhttp_metric_handler_errors_total counter
promhttp_metric_handler_errors_total{cause="encoding"} 0
promhttp_metric_handler_errors_total{cause="gathering"} 0`

		assert.HTTPBodyContains(
			t,
			h,
			http.MethodGet,
			"https://127.0.0.1:50000/",
			nil,
			metricHTTPResponse,
		)
	})
}

func TestMetricCollector_RecordRequestHit(t *testing.T) {
	const testString = "testString"

	t.Parallel()

	t.Run("assert request hit", func(t *testing.T) {
		tests := []struct {
			name       string
			method     string
			path       string
			initialHit float64
			repeat     int
			expected   float64
		}{
			{
				name:       "Record request hit for GET method and / path",
				method:     "GET",
				path:       "/",
				initialHit: 0,
				repeat:     1,
				expected:   1,
			},
			{
				name:       "Record request hit multiple times",
				method:     "POST",
				path:       "/api",
				initialHit: 0,
				repeat:     3,
				expected:   3,
			},
			{
				name:       "Record request hit 1`000 times",
				method:     "POST",
				path:       "/api",
				initialHit: 0,
				repeat:     1000,
				expected:   1000,
			},
			{
				name:       "Record request hit 10`000 times",
				method:     "POST",
				path:       "/api",
				initialHit: 0,
				repeat:     10000,
				expected:   10000,
			},
			{
				name:       "Record request hit 100`000 times",
				method:     "POST",
				path:       "/api",
				initialHit: 0,
				repeat:     100_000,
				expected:   100_000,
			},
			{
				name:       "Record request hit 1`000`000 times",
				method:     "POST",
				path:       "/api",
				initialHit: 0,
				repeat:     1_000_000,
				expected:   1_000_000,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Create a new metrics instance
				m := metrics.NewMetrics(testString)

				// Create a new MetricCollector with a new Prometheus registry
				collector := metrics.NewMetricCollector(
					m,
					prometheus.NewRegistry(),
				)

				// Set initial hits for the path
				m.RequestsHit.WithLabelValues(tt.method, tt.path).Add(tt.initialHit)

				// Call RecordRequestHit to simulate request hits
				for i := 0; i < tt.repeat; i++ {
					collector.RecordRequestHit(tt.method, tt.path)
				}

				// Use prometheus test utility to retrieve the metric and verify its value
				afterHit := testutil.ToFloat64(
					m.RequestsHit.WithLabelValues(tt.method, tt.path))

				assert.Equal(t, tt.expected, afterHit)
			})
		}
	})

	t.Run("assert concurrent request hit", func(t *testing.T) {
		tests := []struct {
			name       string
			method     string
			path       string
			initialHit float64
			repeat     int
			expected   float64
		}{
			{
				name:       "Record request hit for GET method and / path",
				method:     "GET",
				path:       "/",
				initialHit: 0,
				repeat:     1,
				expected:   1,
			},
			{
				name:       "Record request hit multiple times",
				method:     "POST",
				path:       "/api",
				initialHit: 0,
				repeat:     3,
				expected:   3,
			},
			{
				name:       "Record request hit 1`000 times",
				method:     "POST",
				path:       "/api",
				initialHit: 0,
				repeat:     1000,
				expected:   1000,
			},
			{
				name:       "Record request hit 10`000 times",
				method:     "POST",
				path:       "/api",
				initialHit: 0,
				repeat:     10000,
				expected:   10000,
			},
			{
				name:       "Record request hit 100`000 times",
				method:     "POST",
				path:       "/api",
				initialHit: 0,
				repeat:     100_000,
				expected:   100_000,
			},
			{
				name:       "Record request hit 1`000`000 times",
				method:     "POST",
				path:       "/api",
				initialHit: 0,
				repeat:     1_000_000,
				expected:   1_000_000,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Create a new metrics instance
				m := metrics.NewMetrics(testString)

				// Create a new MetricCollector with a new Prometheus registry
				collector := metrics.NewMetricCollector(
					m,
					prometheus.NewRegistry(),
				)

				// Set initial hits for the path
				m.RequestsHit.WithLabelValues(tt.method, tt.path).Add(tt.initialHit)

				var wg sync.WaitGroup

				wg.Add(tt.repeat)

				// Call RecordRequestHit to simulate request hits
				for i := 0; i < tt.repeat; i++ {
					go func() {
						defer wg.Done()
						collector.RecordRequestHit(tt.method, tt.path)
					}()
				}

				wg.Wait()

				// Use prometheus test utility to retrieve the metric and verify its value
				afterHit := testutil.ToFloat64(
					m.RequestsHit.WithLabelValues(tt.method, tt.path))

				assert.Equal(t, tt.expected, afterHit)
			})
		}
	})

	t.Run("assert different uri won`t collide", func(t *testing.T) {
		// Create a new metrics instance
		m := metrics.NewMetrics(testString)

		// Create a new MetricCollector with a new Prometheus registry
		collector := metrics.NewMetricCollector(
			m,
			prometheus.NewRegistry(),
		)

		const path = "/"

		// Set initial hits for the path
		m.RequestsHit.WithLabelValues(http.MethodGet, path).Add(0)

		// Call RecordRequestHit to simulate request hits

		collector.RecordRequestHit(http.MethodGet, path+"somePath")

		// Use prometheus test utility to retrieve the metric and verify its value
		afterHit := testutil.ToFloat64(
			m.RequestsHit.WithLabelValues(http.MethodGet, path))

		assert.Equal(t, float64(0), afterHit)
	})

	t.Run("assert different method won`t collide", func(t *testing.T) {
		// Create a new metrics instance
		m := metrics.NewMetrics(testString)

		// Create a new MetricCollector with a new Prometheus registry
		collector := metrics.NewMetricCollector(
			m,
			prometheus.NewRegistry(),
		)

		const path = "/"

		// Set initial hits for the path
		m.RequestsHit.WithLabelValues(http.MethodGet, path).Add(0)

		// Call RecordRequestHit to simulate request hits

		collector.RecordRequestHit(http.MethodPost, path)

		// Use prometheus test utility to retrieve the metric and verify its value
		afterHit := testutil.ToFloat64(
			m.RequestsHit.WithLabelValues(http.MethodGet, path))

		assert.Equal(t, float64(0), afterHit)
	})
}

func TestMetricCollector_RequestProceedingDuration(t *testing.T) {
	const (
		status       = "200"
		method       = "GET"
		path         = "/"
		testDuration = time.Microsecond
	)

	t.Parallel()

	t.Run("assert request proceeding duration", func(t *testing.T) {
		// Create a new metrics instance
		m := metrics.NewMetrics("test")

		// Create a new MetricCollector with a new Prometheus registry
		collector := metrics.NewMetricCollector(
			m,
			prometheus.NewRegistry(),
		)

		// Call RequestProceedingDuration to record a request duration
		collector.RequestProceedingDuration(status, method, path, testDuration)

		// Use prometheus test utility to retrieve the metric and verify its value
		durMetric, _ := m.RequestProceedingDuration.MetricVec.GetMetricWithLabelValues(status, method, path)
		dto := &io_prometheus_client.Metric{}
		_ = durMetric.Write(dto)

		require.NotZero(t, *dto.Summary.SampleSum, "expected non-zero metric value")

		// Since duration is recorded in seconds, we compare it with the duration in seconds
		assert.Equal(t, testDuration.Seconds(), *dto.Summary.SampleSum, "unexpected metric value")
	})
}
