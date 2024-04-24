package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	mw "service/http/middleware"
)

type Decorator func(router chi.Router) []io.Closer

func addMiddleware(r chi.Router) []io.Closer {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return nil
}

func healthEndpoints(r chi.Router) []io.Closer { //nolint:unparam
	r.Get("/healthz", mw.Healthz)
	r.Get("/readyz", mw.Readyz)
	r.Get("/ping", mw.Ping)

	return nil
}

func addRoutes(r chi.Router) []io.Closer {
	healthEndpoints(r)
	return nil
}
