package middleware

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"service/tracing"
	"testing"
)

func TestMain(m *testing.M) {
	err := tracing.RegisterTracerProvider(context.Background(), resource.Default())
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestResolveTraceIDInHTTP(t *testing.T) {
	t.Run("assert span passed through http", func(t *testing.T) {
		testserver := httptest.NewServer(
			ResolveTraceIDInHTTP("testing")(
				http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
					span := trace.SpanFromContext(r.Context())
					defer span.End()

					ctx := span.SpanContext()
					assert.True(t, ctx.HasTraceID())
					assert.True(t, ctx.IsRemote())
					assert.True(t, ctx.HasSpanID())
				})))
		defer testserver.Close()

		ctx, cancelFunc := context.WithCancel(context.Background())

		resp, e := otelhttp.Get(ctx, testserver.URL)
		require.NoError(t, e)
		defer resp.Body.Close() //nolint:errcheck

		cancelFunc()
	})
	t.Run("assert span generated into middleware", func(t *testing.T) {
		testserver := httptest.NewServer(
			ResolveTraceIDInHTTP("testing")(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				span := trace.SpanFromContext(r.Context())
				defer span.End()

				ctx := span.SpanContext()
				assert.True(t, ctx.HasTraceID())
				assert.False(t, ctx.IsRemote())
				assert.True(t, ctx.HasSpanID())
			})))
		defer testserver.Close()

		resp, e := http.Get(testserver.URL)

		require.NoError(t, e)

		defer resp.Body.Close() //nolint:errcheck
	})
}
