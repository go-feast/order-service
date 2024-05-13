package middleware

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"reflect"
)

func ResolveTraceIDInHTTP(serviceName string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx = r.Context()
			)

			extractedCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
			if reflect.DeepEqual(ctx, extractedCtx) {

				attrs := []attribute.KeyValue{
					semconv.URLFull(r.URL.String()),
				}

				otel.GetTracerProvider().
					Tracer(serviceName).Start(ctx, "http.middleware",
					trace.WithNewRoot(),
					trace.WithSpanKind(trace.SpanKindServer),
					trace.WithAttributes(attrs...),
				)

				otel.GetTextMapPropagator().
					Inject(ctx, propagation.HeaderCarrier(r.Header))
			}

			r = r.WithContext(extractedCtx)

			next(w, r)
		}

		return fn
	}
}
