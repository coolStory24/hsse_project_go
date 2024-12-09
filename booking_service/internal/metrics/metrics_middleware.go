package metrics

import (
	"net/http"
	"time"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		method := r.Method
		path := r.URL.Path

		// ResponseWriter wrapper to capture status code
		ww := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(ww, r)

		duration := time.Since(startTime).Seconds()
		HTTPRequestTotal.WithLabelValues(method, path, http.StatusText(ww.statusCode)).Inc()
		HTTPResponseDuration.WithLabelValues(method, path).Observe(duration)
	})
}

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.statusCode = code
	sw.ResponseWriter.WriteHeader(code)
}
