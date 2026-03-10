package middleware

import (
	"log/slog"
	"net/http"
	// step3 "strconv"
	// step3 "time"
	// step3 "github.com/prometheus/client_golang/prometheus"
	// step6 "go.opentelemetry.io/otel/trace"
)

// step3 // Prometheus metrics — registered once at package init.
// step3 var (
// step3 	httpRequestsTotal = prometheus.NewCounterVec(
// step3 		prometheus.CounterOpts{
// step3 			Name: "http_requests_total",
// step3 			Help: "Total number of HTTP requests by method, path and status code.",
// step3 		},
// step3 		[]string{"method", "path", "status"},
// step3 	)
// step3 	httpRequestDurationSeconds = prometheus.NewHistogramVec(
// step3 		prometheus.HistogramOpts{
// step3 			Name:    "http_request_duration_seconds",
// step3 			Help:    "HTTP request latency distribution.",
// step3 			Buckets: prometheus.DefBuckets,
// step3 		},
// step3 		[]string{"method", "path"},
// step3 	)
// step3 )
// step3
// step3 func init() {
// step3 	prometheus.MustRegister(httpRequestsTotal, httpRequestDurationSeconds)
// step3 }

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

// Metrics is an HTTP middleware that records Prometheus request counters and latency histograms.
// The actual recording is activated at step3.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// step3 start := time.Now()
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		// step3 duration := time.Since(start).Seconds()
		// step3 status := strconv.Itoa(rw.status)
		// step3 httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		// step3 httpRequestDurationSeconds.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

// Logging is an HTTP middleware that emits a structured JSON log line for every request.
// The actual log call is activated at step4; trace_id field is added at step6.
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			// step4 traceID := ""
			// step6 if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
			// step6 	traceID = span.SpanContext().TraceID().String()
			// step6 }
			// step4 logger.InfoContext(r.Context(), "http request",
			// step4 	"method", r.Method,
			// step4 	"path", r.URL.Path,
			// step4 	"status", rw.status,
			// step4 	"trace_id", traceID,
			// step4 )
		})
	}
}

// ensure slog is used at compile time even before step4 is activated
var _ = slog.Default
