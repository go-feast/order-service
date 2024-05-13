// Package metrics provide functions to create metrics in global prometheus registry
// Before initializing any of metrics you should call RegisterServiceName.
//
//	metrics.RegisterServiceName(serviceName)
//	...
//	counter = metrics.NewCounter(pkgName, "counter_total")
//	counter.Inc()
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

// serviceName states for service application name. Used in all metrics as Namespace.
var serviceName string

// RegisterServiceName registers application service name for all metrics.
// Must be called before initializing any of metrics.
func RegisterServiceName(name string) {
	serviceName = name
}

// NewCounter register Counter in global registry.
func NewCounter(subsystem, name string) prometheus.Counter {
	return promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: serviceName,
			Subsystem: subsystem,
			Name:      name,
		},
	)
}

// NewCounterVec register CounterVec in global registry.
func NewCounterVec(subsystem, name string, labels ...string) *prometheus.CounterVec {
	return promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: serviceName,
			Subsystem: subsystem,
			Name:      name,
		},
		labels,
	)
}

// NewGauge register Gauge in global registry.
func NewGauge(subsystem, name string) prometheus.Gauge {
	return promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: serviceName,
			Subsystem: subsystem,
			Name:      name,
		},
	)
}

// NewGaugeVec register GaugeVec in global registry.
func NewGaugeVec(subsystem, name string, labels ...string) *prometheus.GaugeVec {
	return promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: serviceName,
			Subsystem: subsystem,
			Name:      name,
		},
		labels,
	)
}

// NewSummary register Summary in global registry.
func NewSummary(subsystem, name string, buckets map[float64]float64, maxAge time.Duration) prometheus.Summary {
	return promauto.NewSummary(
		prometheus.SummaryOpts{
			Namespace:  serviceName,
			Subsystem:  subsystem,
			Name:       name,
			Objectives: buckets,
			MaxAge:     maxAge,
		},
	)
}

// NewSummaryVec register SummaryVec in global registry.
func NewSummaryVec(subsystem, name string, labels ...string) *prometheus.SummaryVec {
	return promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: serviceName,
			Subsystem: subsystem,
			Name:      name,
		},
		labels,
	)
}

// NewHistogram register Histogram in global registry.
func NewHistogram(subsystem, name string, buckets []float64) prometheus.Histogram {
	return promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: serviceName,
			Subsystem: subsystem,
			Name:      name,
			Buckets:   buckets,
		},
	)
}

// NewHistogramVec register HistogramVec in global registry.
func NewHistogramVec(subsystem, name string, buckets []float64, labels ...string) *prometheus.HistogramVec {
	return promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: serviceName,
			Subsystem: subsystem,
			Name:      name,
			Buckets:   buckets,
		},
		labels,
	)
}
