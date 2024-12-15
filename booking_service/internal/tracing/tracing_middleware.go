package tracing

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("http-server")
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		spanCtx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path)
		defer span.End()

		if sc := span.SpanContext(); sc.IsValid() {
			w.Header().Set("X-Trace-ID", sc.TraceID().String())
		}

		next.ServeHTTP(w, r.WithContext(spanCtx))
	})
}
