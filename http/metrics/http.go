// Package metrics provides tools and middlewares
// for manipulating metrics around Prometheus
// TODO: If need - possible to make otlp refactoring.
package metrics

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/prometheus/client_golang/prometheus/promauto" //nolint:revive
	"net"
	"service/config"

	//nolint:revive
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	server        *http.Server
	metricService *MetricService
	l             *zap.Logger
}

var metricServer *Server

var once sync.Once

func NewMetricsServer(serviceName string, c *config.MetricServerConfig, l *zap.Logger) *Server {
	once.Do(func() {
		var (
			metricService = NewMetricsService(NewMetrics(serviceName))
		)

		router := chi.NewRouter()

		router.Get("/metrics", metricService.Handler())

		metricServer = &Server{
			server: &http.Server{ //nolint:gosec
				Addr:         net.JoinHostPort(c.Host, c.Port),
				Handler:      router,
				WriteTimeout: c.WriteTimeout,
				ReadTimeout:  c.ReadTimeout,
				IdleTimeout:  c.IdleTimeout,
			},
			metricService: metricService,
			l:             l,
		}
	})

	return metricServer
}

// ListenAndServe non-blocking server function for listening on requests.
// Server shuts down when ctx is done.
func (s *Server) ListenAndServe(ctx context.Context) {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		s.l.Info("metric server is running", zap.String("url", "http://"+s.server.Addr))

		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return
		}

		s.l.Info("metric server is shutting down...", zap.Error(err))
	}()

	go func() {
		defer wg.Done()

		<-ctx.Done()

		err := s.server.Shutdown(ctx)
		if err != nil &&
			!errors.Is(err, http.ErrServerClosed) &&
			!errors.Is(err, context.Canceled) {
			s.l.Error("metric server: shutdown error:", zap.Error(err))
			return
		}

		wg.Wait()
	}()
}

func MetricRequests(next http.Handler) http.Handler {
	// rps
	// response time
	// status metric
	fn := func(w http.ResponseWriter, r *http.Request) {
		response := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		beforeRequest := time.Now()

		defer func() {
			afterRequest := time.Since(beforeRequest)
			status := http.StatusText(response.Status())
			metricServer.metricService.RequestProceedingDuration(status, r.Method, afterRequest)
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
