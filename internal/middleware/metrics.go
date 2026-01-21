package middleware

import (
	"avito/internal/logger"
	"net/http"
	"time"

	"avito/internal/metrics"

	chi "github.com/go-chi/chi/v5"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsAndLogsMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}
			next.ServeHTTP(rw, r)
			duration := time.Since(start)

			metrics.HttpRequestsTotal.WithLabelValues(
				r.Method,
				chi.RouteContext(r.Context()).RoutePattern(),
				http.StatusText(rw.status),
			).Inc()

			metrics.HttpRequestDuration.WithLabelValues(
				chi.RouteContext(r.Context()).RoutePattern(),
			).Observe(duration.Seconds())

			log.Info("http request",
				logger.Field{Key: "method", Value: r.Method},
				logger.Field{Key: "path", Value: r.URL.Path},
				logger.Field{Key: "status", Value: rw.status},
				logger.Field{Key: "duration", Value: duration},
			)
		})
	}
}
