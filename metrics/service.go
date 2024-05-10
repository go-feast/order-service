package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"reflect"
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
	ms.reg.MustRegister(ms.Collectors()...)

	return promhttp.HandlerFor(ms.reg, promhttp.HandlerOpts{Registry: ms.reg}).ServeHTTP
}

func (ms *MetricCollector) RequestProceedingDuration(status, method, uri string, dur time.Duration) {
	ms.m.RequestProceedingDuration.WithLabelValues(status, method, uri).Observe(dur.Seconds())
}

func (ms *MetricCollector) RecordRequestHit(method, uri string) {
	ms.m.RequestsHit.WithLabelValues(method, uri).Inc()
}

func (ms *MetricCollector) Collectors() []prometheus.Collector {
	if ms.m == nil {
		return nil
	}

	v := reflect.ValueOf(*(ms.m))

	// getting number of fields
	n := v.NumField()

	collectors := make([]prometheus.Collector, n)

	for i := 0; i < n; i++ {
		field := v.Field(i)

		collector, ok := field.Interface().(prometheus.Collector)
		if ok {
			collectors[i] = collector
		}
	}

	return collectors
}
