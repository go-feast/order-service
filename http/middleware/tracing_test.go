package middleware

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"os"
	"service/config"
	serv "service/server"
	"service/tracing"
	"testing"
)

func TestMain(m *testing.M) {
	err := tracing.RegisterTracerProvider(context.Background(), "test")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestResolveTraceIDInHTTP(t *testing.T) {
	t.Run("assert span passed through http", func(t *testing.T) {
		ctx := context.Background()

		server, router := serv.NewServer(&config.MainServiceServerConfig{Host: "127.0.0.1", Port: "40000"})

		router.
			With(ResolveTraceIDInHTTP("testing")).
			Get("/", func(_ http.ResponseWriter, r *http.Request) {
				span := trace.SpanFromContext(r.Context())
				defer span.End()

				ctx := span.SpanContext()
				assert.True(t, ctx.HasTraceID())
				assert.True(t, ctx.IsRemote())
				assert.True(t, ctx.HasSpanID())
			})

		ctx, cancelFunc := context.WithCancel(ctx)
		started, _ := serv.Run(ctx, server)

		<-started

		resp, err := otelhttp.Get(ctx, "http://127.0.0.1:40000/")
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck

		cancelFunc()
	})
	t.Run("assert span generated into middleware", func(t *testing.T) {
		ctx := context.Background()

		server, router := serv.NewServer(&config.MainServiceServerConfig{Host: "127.0.0.1", Port: "40000"})

		router.
			With(ResolveTraceIDInHTTP("testing")).
			Get("/", func(_ http.ResponseWriter, r *http.Request) {
				span := trace.SpanFromContext(r.Context())
				defer span.End()

				ctx := span.SpanContext()
				assert.True(t, ctx.HasTraceID())
				assert.False(t, ctx.IsRemote())
				assert.True(t, ctx.HasSpanID())
			})

		ctx, cancelFunc := context.WithCancel(ctx)
		started, _ := serv.Run(ctx, server)

		<-started

		resp, err := http.Get("http://127.0.0.1:40000/")
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck

		cancelFunc()
	})
}
