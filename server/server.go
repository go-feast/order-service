package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
	"net/http"
	"service/config"
	"service/logging"
	"sync"
)

// Run function run all servers that we provide and handles graceful shutdown via context.
func Run(ctx context.Context, servers ...*http.Server) (started chan struct{}, err chan error) {
	started, err = make(chan struct{}), make(chan error, len(servers)<<1)

	var l = logging.New()

	go func() {
		errIn := make(chan error)

		// collecting errors
		go func() {
			var ex error
			for e := range errIn {
				ex = errors.Join(ex, e)
			}

			err <- ex

			close(err)
		}()

		group, _ := errgroup.WithContext(ctx)

		var wg sync.WaitGroup

		wg.Add(len(servers))

		for _, server := range servers {
			group.Go(func() error {
				wg.Done()

				err := server.ListenAndServe()
				if err != nil {
					errIn <- err
				}

				return err
			})
		}

		go func() {
			<-ctx.Done()

			l.Info().Msg("shutting down servers")

			for _, server := range servers {
				err := server.Shutdown(ctx)
				if err != nil {
					errIn <- err
				}
			}

			err := group.Wait()
			if err != nil {
				errIn <- err
			}

			close(errIn)

			l.Info().Msg("servers shut down")
		}()

		wg.Wait()

		l.Info().Strs("urls", getAddrs(servers...)).Msg("running servers")

		close(started)
	}()

	return started, err
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

func getAddrs(s ...*http.Server) []string {
	addrs := make([]string, len(s))
	for i, server := range s {
		addrs[i] = "http://" + server.Addr
	}

	return addrs
}
