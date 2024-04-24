package middleware_test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"service/http/middleware"
	"testing"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()
	t.Run("healthz", testRoute(t, middleware.Healthz))
	t.Run("readyz", testRoute(t, middleware.Readyz))
	t.Run("ping", testRoute(t, middleware.Ping))
}

func testRoute(t *testing.T, h http.HandlerFunc) func(t *testing.T) {
	return func(_ *testing.T) {
		ts := httptest.NewServer(h)
		defer ts.Close()

		assert.HTTPStatusCode(t, h,
			"GET", "http://localhost:8080/", nil, 200)
	}
}
