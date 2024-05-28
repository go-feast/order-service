package middleware

import (
	"fmt"
	"net/http"
	"service/metrics"
	"strconv"
	"time"
)

const (
	subsystem = "http"
	method    = "method"
	uri       = "uri"
	code      = "code"
)

type wrappedResponseWriter struct {
	w       http.ResponseWriter
	code    string
	codeInt int
}

func (w *wrappedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *wrappedResponseWriter) Write(bytes []byte) (int, error) {
	return w.w.Write(bytes)
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.codeInt = statusCode
	w.code = strconv.Itoa(statusCode)
	w.w.WriteHeader(statusCode)
}

func RecordRequestHit(handlerName string) func(http.Handler) http.Handler {
	metric := fmt.Sprintf("%s_request_hit_total", handlerName)

	var recorder = metrics.NewCounterVec(
		subsystem, metric,
		method, uri, code,
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var wrapper = &wrappedResponseWriter{w: w}
			defer func() {
				recorder.WithLabelValues(r.Method, r.RequestURI, wrapper.code).Inc()
			}()
			next.ServeHTTP(wrapper, r)
		})
	}
}

func RecordRequestDuration(handlerName string) func(http.Handler) http.Handler {
	metric := fmt.Sprintf("%s_request_duration", handlerName)
	recorder := metrics.NewHistogramVec(
		subsystem, metric,
		[]float64{0.1, 0.25, 0.5, 0.75, 1},
		method, uri, code,
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			beforeRequest, wrapper := time.Now(), &wrappedResponseWriter{w: w}
			defer func() {
				afterRequest := time.Since(beforeRequest)
				recorder.WithLabelValues(r.Method, r.RequestURI, wrapper.code).
					Observe(afterRequest.Seconds())
			}()
			next.ServeHTTP(wrapper, r)
		})
	}
}
