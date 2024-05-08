package server_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"service/http/middleware"
	"service/logging"
	serv "service/server"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("graceful shutdown of the servers", func(t *testing.T) {
		c1 := &testConfig{hostPort: "127.0.0.1:50000"}
		s1, r1 := serv.NewServer(c1)
		c2 := &testConfig{hostPort: "127.0.0.1:50001"}
		s2, r2 := serv.NewServer(c2)

		ctx, cancel := context.WithCancel(context.Background())

		l, _ := logging.NewLogger()

		started, err := serv.Run(ctx, l, s1, s2)

		<-started

		r1.Get("/", middleware.Healthz)
		r2.Get("/", middleware.Healthz)

		resp1, _ := http.Get(c1.URL())
		resp1.Body.Close() //nolint:errcheck

		resp2, _ := http.Get(c2.URL())
		resp2.Body.Close() //nolint:errcheck

		assert.Equal(t, http.StatusOK, resp1.StatusCode)
		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		cancel()

		e := <-err

		assert.ErrorIs(t, e, http.ErrServerClosed)
	})
}

func TestNewServer(t *testing.T) {
	t.Run("assert adding routers affect on server", func(t *testing.T) {
		serverConfig := &testConfig{hostPort: "127.0.0.1:50000"}

		server, router := serv.NewServer(serverConfig)

		router.Get("/", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

		assert.HTTPSuccess(t, server.Handler.ServeHTTP, http.MethodGet, serverConfig.URL(), nil)
	})

	t.Run("panic when invalid config passed", func(t *testing.T) {
		serverConfig := &panicConfig{}

		assert.Panics(t, func() { serv.NewServer(serverConfig) })
	})
}

type panicConfig struct{}

func (p panicConfig) HostPort() string {
	panic("implement me")
}

func (p panicConfig) WriteTimeoutDur() time.Duration {
	panic("implement me")
}

func (p panicConfig) ReadTimeoutDur() time.Duration {
	panic("implement me")
}

func (p panicConfig) IdleTimeoutDur() time.Duration {
	panic("implement me")
}

func (p panicConfig) ReadHeaderTimeoutDur() time.Duration {
	panic("implement me")
}

type testConfig struct {
	hostPort string
}

func (t *testConfig) URL() string {
	return fmt.Sprintf("http://%s/", t.HostPort())
}

func (t *testConfig) HostPort() string {
	return t.hostPort
}

func (t *testConfig) WriteTimeoutDur() time.Duration {
	return 0
}

func (t *testConfig) ReadTimeoutDur() time.Duration {
	return 0
}

func (t *testConfig) IdleTimeoutDur() time.Duration {
	return 0
}

func (t *testConfig) ReadHeaderTimeoutDur() time.Duration {
	return 0
}
