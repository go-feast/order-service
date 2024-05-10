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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	// TODO: make atomic.Value
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

	tracer := otel.Tracer("metric tracer")
	var h = metricServer.metricService.Handler()
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
		)

		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

		ctx, span := tracer.Start(ctx, "metric handler") //nolint:ineffassign

		defer span.End()

		h(w, r)
	}
}

func RecordRequestHit(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		metricServer.metricService.RecordRequestHit(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
