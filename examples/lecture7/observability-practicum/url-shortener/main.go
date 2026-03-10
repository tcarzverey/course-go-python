package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	// step2 "github.com/prometheus/client_golang/prometheus/promhttp"
	// step5 "github.com/observability-practicum/url-shortener/internal/telemetry"
	// step6 "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/observability-practicum/url-shortener/internal/handler"
	"github.com/observability-practicum/url-shortener/internal/middleware"
	"github.com/observability-practicum/url-shortener/internal/storage"
)

func main() {
	// ── Configuration ────────────────────────────────────────────────────────
	port := envOr("PORT", "8080")
	baseURL := envOr("BASE_URL", "http://localhost:"+port)
	// step5 otlpEndpoint := envOr("OTEL_EXPORTER_OTLP_ENDPOINT", "otel-collector:4317")
	// step7 statsServiceURL := envOr("STATS_SERVICE_URL", "http://stats-service:8081")

	// ── Logger ───────────────────────────────────────────────────────────────
	// Base (step0-3): structured logger disabled, use slog.Default() (plain text to stdout)
	logger := slog.Default()

	// step4 // Replace the logger above with a JSON handler so Loki can parse fields.
	// step4 logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	// step4 slog.SetDefault(logger)

	// ── OpenTelemetry tracer ─────────────────────────────────────────────────
	// step5 ctx := context.Background()
	// step5 shutdown, err := telemetry.InitTracer(ctx, "url-shortener", otlpEndpoint)
	// step5 if err != nil {
	// step5 	log.Fatalf("failed to initialise tracer: %v", err)
	// step5 }
	// step5 defer func() { _ = shutdown(context.Background()) }()

	// ── Storage & Handler ─────────────────────────────────────────────────────
	store := storage.New()

	// step7 // Pass the stats-service URL so the handler can track clicks there.
	// step7 statsURL := statsServiceURL
	statsURL := "" // step7: replace this line with the two lines above

	h := handler.New(store, baseURL, logger, statsURL)

	// ── Router ────────────────────────────────────────────────────────────────
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("GET /stats/{code}", h.Stats) // must be before /{code}
	mux.HandleFunc("GET /{code}", h.Redirect)
	mux.HandleFunc("GET /health", h.Health)
	// step2 mux.Handle("GET /metrics", promhttp.Handler())

	// ── Middleware chain (outermost runs first) ───────────────────────────────
	var root http.Handler = mux

	// step6 // OTel HTTP middleware: creates a span for every request automatically.
	// step6 root = otelhttp.NewHandler(root, "url-shortener-http",
	// step6 	otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	// step6 )

	root = middleware.Metrics(root)              // step3 activates recording
	root = middleware.Logging(logger)(root)       // step4 activates log output

	// ── Start server ──────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", port)
	log.Printf("starting url-shortener on %s (base_url=%s)", addr, baseURL)
	if err := http.ListenAndServe(addr, root); err != nil {
		log.Fatalf("server error: %v", err)
	}

	// keep context import used before step5
	_ = context.Background
	_ = os.Stderr
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
