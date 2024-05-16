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

func ResolveTraceIDInHTTP(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx = r.Context()
			)
			// TODO: maybe wrap status code
			wrapped := &wrappedResponseWriter{w: w}

			extractedCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
			if reflect.DeepEqual(ctx, extractedCtx) {
				attrs := []attribute.KeyValue{
					semconv.HTTPMethod(r.Method),
					semconv.HTTPScheme(r.URL.Scheme),
					semconv.HTTPUserAgent(r.UserAgent()),
					semconv.URLPath(r.URL.String()),
				}

				var span trace.Span

				extractedCtx, span = otel.GetTracerProvider().
					Tracer(serviceName).
					Start(ctx, "http.middleware",
						trace.WithNewRoot(),
						trace.WithSpanKind(trace.SpanKindServer),
						trace.WithAttributes(attrs...),
					)
				defer span.End()

				otel.GetTextMapPropagator().
					Inject(extractedCtx, propagation.HeaderCarrier(r.Header))
			}

			r = r.WithContext(extractedCtx)

			next.ServeHTTP(wrapped, r)
		})

		return fn
	}
}
