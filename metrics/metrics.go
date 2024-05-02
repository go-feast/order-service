package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
)

// Metrics struct contains metrics fields
// For operations(increasing, decreasing etc.) on them - responsible MetricService.
// Metrics should only contain prometheus types of metrics
type Metrics struct {
	// TODO: ask
	RequestProceedingDuration *prometheus.CounterVec

	// TODO: ask
	// Should be used for getting number of request by rate()
	RequestsHit *prometheus.CounterVec
}

func NewMetrics(serviceName string) *Metrics {
	return &Metrics{
		RequestProceedingDuration: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: serviceName,
			Name:      "request_processing_duration",
			Help:      "Metric represents duration of proceeding request in nanoseconds. Portioned by status, method, uri. ",
		}, []string{"status", "method", "uri"}),
		RequestsHit: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: serviceName,
			Name:      "request_hit_total",
			Help:      "Metric represents hits for specified request. Portioned by method, uri. ",
		}, []string{"method", "uri"}),
	}
}

func (m *Metrics) Collectors() []prometheus.Collector {
	if m == nil {
		return nil
	}

	v := reflect.ValueOf(*m)

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
