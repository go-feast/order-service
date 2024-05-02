// Package httpmetrics provides tools and middlewares
// for manipulating metrics around Prometheus
// TODO: If need - possible to make otlp refactoring.
package httpmetrics

import (
	_ "github.com/prometheus/client_golang/prometheus/promauto" //nolint:revive
	"go.uber.org/zap"
	"net/http"
	"service/metrics"
)

type Server struct {
	server        *http.Server
	metricService *metrics.MetricService
	l             *zap.Logger
}

func RegisterServer(server *http.Server, metricService *metrics.MetricService, l *zap.Logger) {
	metricServer = &Server{server: server, metricService: metricService, l: l}
}

var (
	metricServer *Server
)

func RecordRequestHit(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		metricServer.metricService.RecordRequestHit(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
