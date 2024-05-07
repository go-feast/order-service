package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

// MetricCollector responsible for collecting metrics via methods and expose a http.HandlerFunc.
type MetricCollector struct {
	m   *Metrics
	reg *prometheus.Registry
}

func NewMetricCollector(metrics *Metrics, reg *prometheus.Registry) *MetricCollector {
	return &MetricCollector{m: metrics, reg: reg}
}

func (ms *MetricCollector) Handler() http.HandlerFunc {
	ms.reg.MustRegister(ms.m.Collectors()...)

	return promhttp.HandlerFor(ms.reg, promhttp.HandlerOpts{Registry: ms.reg}).ServeHTTP
}

func (ms *MetricCollector) RequestProceedingDuration(status, method, uri string, dur time.Duration) {
	ms.m.RequestProceedingDuration.WithLabelValues(status, method, uri).Observe(dur.Seconds())
}

func (ms *MetricCollector) RecordRequestHit(method, uri string) {
	ms.m.RequestsHit.WithLabelValues(method, uri).Inc()
}
