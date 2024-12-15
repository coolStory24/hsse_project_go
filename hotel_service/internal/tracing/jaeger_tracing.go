package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.21.0"
)

func InitTracerProvider(serviceName, jaegerEndpoint string) (*trace.TracerProvider, error) {
	// Create a Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return nil, err
	}

	// Create a trace provider with the exporter
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Set the global tracer provider
	otel.SetTracerProvider(tp)
	return tp, nil
}

func ShutdownTracerProvider(ctx context.Context, tp *trace.TracerProvider) {
	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown tracer provider: %v", err)
	}
}
