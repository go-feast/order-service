package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"service/config"
	"sync"
)

// Run function run all servers that we provide and handles graceful shutdown via context.
func Run(ctx context.Context, l *zap.Logger, servers ...*http.Server) (err error) {
	var wg sync.WaitGroup

	wg.Add(len(servers))
	addrs := make([]string, len(servers))

	for _, server := range servers {
		addrs = append(addrs, "http://"+server.Addr)
	}

	for _, server := range servers {
		go func(s *http.Server) {
			defer wg.Done()

			err = s.ListenAndServe()
		}(server)
	}

	l.Info("running servers", zap.Strings("urls", addrs))

	wg.Add(1)
	// function for gracefully shutdown
	go func(ctx context.Context, servers ...*http.Server) {
		defer wg.Done()
		<-ctx.Done()

		l.Info("shutting down servers")

		for _, server := range servers {
			err = server.Shutdown(ctx)
		}
	}(ctx, servers...)

	wg.Wait()

	l.Info("servers shut down")

	return
}

func NewServer(c config.ServerConfig) (*http.Server, chi.Router) {
	mux := chi.NewRouter()

	server := &http.Server{
		Addr:              c.HostPort(),
		Handler:           mux,
		ReadTimeout:       c.ReadTimeoutDur(),
		ReadHeaderTimeout: c.ReadHeaderTimeoutDur(),
		WriteTimeout:      c.WriteTimeoutDur(),
		IdleTimeout:       c.IdleTimeoutDur(),
	}

	return server, mux
}
