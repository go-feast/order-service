package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type MetricService struct {
	m   *Metrics
	reg *prometheus.Registry
}

func NewMetricsService(metrics *Metrics, reg *prometheus.Registry) *MetricService {
	return &MetricService{m: metrics, reg: reg}
}

func (ms *MetricService) Handler() http.HandlerFunc {
	ms.reg.MustRegister(ms.m.Collectors()...)

	return promhttp.HandlerFor(ms.reg, promhttp.HandlerOpts{Registry: ms.reg}).ServeHTTP
}

func (ms *MetricService) RequestProceedingDuration(_, _ string, _ time.Duration) {
	panic("not implemented")
}

func (ms *MetricService) RecordRequestHit(method, uri string) {
	ms.m.RequestsHit.WithLabelValues(method, uri).Inc()
}
