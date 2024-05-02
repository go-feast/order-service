package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type Metrics struct {
	RequestProceedingDuration *prometheus.GaugeVec
}

func (m *Metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.RequestProceedingDuration,
	}
}

func NewMetrics(serviceName string) *Metrics {
	return &Metrics{
		RequestProceedingDuration: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: serviceName + "_response_duration",
			Help: "Metric represents duration of proceeding request in nanoseconds. Portioned by status, method, uri. ",
		}, []string{"status", "method"}),
	}
}

type MetricService struct {
	m *Metrics
}

func (ms *MetricService) Handler() http.HandlerFunc {
	reg := prometheus.NewRegistry()
	reg.MustRegister(ms.m.Collectors()...)

	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP
}

func (ms *MetricService) RequestProceedingDuration(status, meth string, dur time.Duration) {
	ms.m.RequestProceedingDuration.WithLabelValues(status, meth).Set(float64(dur))
}

func NewMetricsService(metrics *Metrics) *MetricService {
	return &MetricService{m: metrics}
}
