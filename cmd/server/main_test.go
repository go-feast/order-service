package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"service/config"
	"service/logging"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	t.Parallel()

	c := &config.Config{ServerConfig: &config.ServerConfig{Host: "localhost", Port: "8080"}}
	l, _ := logging.NewLogger(config.Development)
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error)

	go func() {
		errCh <- app(ctx, c, l)
	}()

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8080/healthz")
	defer func() {
		e := resp.Body.Close()
		if e != nil {
			t.Error(e)
		}
	}()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cancel()
	assert.NoError(t, <-errCh)
}

func TestDecorators(t *testing.T) {
	r := chi.NewRouter()
	closers := decorators(r, addRoutes)

	assert.NotNil(t, closers)
}

func TestHealthEndpoints(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	healthEndpoints(r)

	urls := [...]string{"/healthz", "/readyz", "/ping"}

	for _, url := range urls {
		assert.HTTPStatusCode(t, r.ServeHTTP, "GET",
			url, nil, 200)
	}
}
