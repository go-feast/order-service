package middleware

import "net/http"

var f http.HandlerFunc = func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

var (
	Healthz = f
	Readyz  = f
	Ping    = f
)
