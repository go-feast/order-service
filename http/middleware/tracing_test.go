package middleware

import (
	"context"
	"log"
	"os"
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
