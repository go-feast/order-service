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
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"io"
	"os"
	"service/config"
)

// RegisterTracerProvider registers global trace.TracerProvider
func RegisterTracerProvider(ctx context.Context, serviceName string) error {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
		resource.WithProcess(),
		resource.WithOS(),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

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
		return newStdOutExporter(os.Stdout)
	case config.Development, config.Local, config.Production:
		traceExporter, err := otlptracehttp.New(ctx,
			// OTEL_EXPORTER_OTLP_TRACES_ENDPOINT MUST be set
			otlptracehttp.WithInsecure(), // HTTPS -> HTTP
		)
		if err != nil {
			return nil, err
		}

		return traceExporter, nil
	}

	return nil, fmt.Errorf("failed to create exporter")
}

func newStdOutExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint(),
	)
}
