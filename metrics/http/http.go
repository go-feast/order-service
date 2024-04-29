// Package http provides tools and middlewares
// for manipulating metrics around Prometheus
// TODO: If need - possible to make otlp refactoring.
package http

import (
	_ "github.com/prometheus/client_golang/prometheus/promauto" //nolint:revive
	_ "github.com/prometheus/client_golang/prometheus/promhttp" //nolint:revive
)
