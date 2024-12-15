package tracing

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitTracerProvider(t *testing.T) {
	tracer, err := InitTracerProvider("hotel-service", "http://localhost:1234/jaeger")

	assert.NoError(t, err)
	assert.NotNil(t, tracer)

	ShutdownTracerProvider(context.Background(), tracer)
}
