package metrics

import "go.opentelemetry.io/otel/sdk/metric"

// otlp metrics
func RegisterMetricExporter() error {

	metric.NewMeterProvider(metric.WithReader())
	return nil
}
