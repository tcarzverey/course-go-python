package middleware

import (
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging is an HTTP middleware that emits a structured log line for every request.
// The actual log call is activated at step5; trace_id field is added at step6.
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			traceID := ""
			if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
				traceID = span.SpanContext().TraceID().String()
			}
			logger.InfoContext(r.Context(), "http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"trace_id", traceID,
			)
		})
	}
}

// ensure slog is used at compile time even before step5 is activated
var _ = slog.Default
