// Package httpmetrics represents singleton pattern of metrics server,
// which exposes handler for reading metrics via http
// and manipulating them in application.
// Example:
//
//	s := &http.Server{...}
//	m := metrics.NewMetrics("serviceName")
//	collector := metrics.NewCollector(m)
//
//	httpmetrics.RegisterServer(s, collector, l)
//	...
//	mux.Get("/metrics", httpmetrics.Handler())
package httpmetrics

// TODO: If need - possible to make otlp refactoring.
import (
	"net/http"
	"service/logging"
	"service/metrics"
)

// Server represents metrics server struct encapsulating metrics.MetricCollector.
// All functions and methods in package should relay on singleton instance of Server.
type Server struct {
	metricService *metrics.MetricCollector
	l             *logging.Logger
}

var (
	metricServer *Server
)

func RegisterServer(metricService *metrics.MetricCollector, l *logging.Logger) {
	metricServer = &Server{metricService: metricService, l: l}
}

func UnregisterServer() {
	metricServer = nil
}

func Handler() http.HandlerFunc {
	if metricServer == nil {
		panic("metric Server not registered")
	}

	return metricServer.metricService.Handler()
}

func RecordRequestHit(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		metricServer.metricService.RecordRequestHit(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
