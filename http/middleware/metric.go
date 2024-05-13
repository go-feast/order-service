package middleware

import (
	"fmt"
	"net/http"
	"service/metrics"
	"time"
)

func RecordRequestHit(handlerName string) func(http.Handler) http.Handler {
	metric := fmt.Sprintf("%s_request_hit_total", handlerName)

	var recorder = metrics.NewCounterVec(
		"http", metric, "method", "url")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder.WithLabelValues(r.Method, r.RequestURI).Inc()
			next.ServeHTTP(w, r)
		})
	}
}

func RecordRequestDuration(handlerName string) func(http.Handler) http.Handler {
	metric := fmt.Sprintf("%s_request_duration", handlerName)
	recorder := metrics.NewHistogramVec(
		"http", metric,
		[]float64{0.1, 0.25, 0.5, 0.75, 1},
		"code", "method",
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			beforeRequest := time.Now()
			defer func() {
				afterRequest := time.Since(beforeRequest)
				recorder.WithLabelValues(r.Method, r.RequestURI).Observe(afterRequest.Seconds())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
