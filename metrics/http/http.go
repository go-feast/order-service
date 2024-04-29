// Package http provides tools and middlewares
// for manipulating metrics around Prometheus
// TODO: If need - possible to make otlp refactoring.
package http

import (
	_ "github.com/prometheus/client_golang/prometheus/promauto"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
)
