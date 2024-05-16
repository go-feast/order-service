package tracing

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"os"
	"service/config"
)

// RegisterTracerProvider registers global trace.TracerProvider
func RegisterTracerProvider(ctx context.Context, res *resource.Resource) error {
	exporter, err := newExporter(ctx)
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		),
	)

	return nil
}

func newExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	switch config.MustGetEnvironment() {
	case config.Testing:
		return stdouttrace.New(stdouttrace.WithWriter(os.Stdout), stdouttrace.WithPrettyPrint())
	case config.Development, config.Local, config.Production:
		// OTEL_EXPORTER_OTLP_TRACES_ENDPOINT MUST be set
		traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithInsecure()) // HTTPS -> HTTP
		if err != nil {
			return nil, err
		}

		return traceExporter, nil
	}

	return nil, fmt.Errorf("failed to create exporter")
}
