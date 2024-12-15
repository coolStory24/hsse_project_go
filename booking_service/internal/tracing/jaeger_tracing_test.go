package tracing

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitTracerProvider(t *testing.T) {
	tracer, err := InitTracerProvider("booking_service", "http://localhost:1234/something")

	assert.NoError(t, err)
	assert.NotNil(t, tracer)

	ShutdownTracerProvider(context.Background(), tracer)
}

func TestTracingMiddleware(t *testing.T) {
	nextHandler := TracingMiddleware(http.NotFoundHandler())
	nextHandler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	assert.NotNil(t, nextHandler)
}
