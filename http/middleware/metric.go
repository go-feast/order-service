package middleware

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"service/metrics"
)

func RecordRequestHit(next http.HandlerFunc) http.HandlerFunc {
	var recorder = metrics.NewCounterVec(
		"http", "request_hit_total", "code", "method")
	return promhttp.InstrumentHandlerCounter(recorder, next)
}

func RecordRequestDuration(next http.HandlerFunc) http.HandlerFunc {
	recorder := metrics.NewHistogramVec(
		"http", "request_duration",
		[]float64{0.1, 0.25, 0.5, 0.75, 1},
		"code", "method",
	)

	return promhttp.InstrumentHandlerDuration(recorder, next)
}
