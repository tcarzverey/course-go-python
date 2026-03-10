package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/observability-practicum/url-shortener/internal/telemetry"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/observability-practicum/url-shortener/internal/handler"
	"github.com/observability-practicum/url-shortener/internal/middleware"
	"github.com/observability-practicum/url-shortener/internal/storage"
)

func main() {
	// ── Configuration ────────────────────────────────────────────────────────
	port := envOr("PORT", "8080")
	baseURL := envOr("BASE_URL", "http://localhost:"+port)
	otlpEndpoint := envOr("OTEL_EXPORTER_OTLP_ENDPOINT", "otel-collector:4317")
	// step7 statsServiceURL := envOr("STATS_SERVICE_URL", "http://stats-service:8081")

	// ── Logger ───────────────────────────────────────────────────────────────
	// Base (step0-1): plain default logger, no structured output.
	logger := slog.Default()

	// Switch to OTel-backed slog — logs are sent via OTLP once step2 is active.
	logger = slog.New(otelslog.NewHandler("url-shortener"))
	slog.SetDefault(logger)

	// ── OpenTelemetry SDK ────────────────────────────────────────────────────
	ctx := context.Background()
	shutdown, err := telemetry.Init(ctx, "url-shortener", otlpEndpoint)
	if err != nil {
		log.Fatalf("failed to initialise OTel SDK: %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	// ── Storage & Handler ─────────────────────────────────────────────────────
	// step8 dbDSN := envOr("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/shortener?sslmode=disable")

	var store storage.Store = storage.New() // in-memory, steps 0–7

	// step8 pgStore, err := storage.NewPGX(ctx, dbDSN)
	// step8 if err != nil { log.Fatalf("db: %v", err) }
	// step8 store = pgStore

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

	// ── Middleware chain (outermost runs first) ───────────────────────────────
	var root http.Handler = mux

	root = middleware.Logging(logger)(root) // step5 activates log output

	// OTel HTTP middleware: outermost — creates span and injects it into context
	// before Logging runs, so trace_id is available in r.Context() inside Logging.
	root = otelhttp.NewHandler(root, "url-shortener-http",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)

	// ── Start server ──────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", port)
	slog.Info("starting url-shortener", "addr", addr, "baseURL", baseURL)
	if err := http.ListenAndServe(addr, root); err != nil {
		log.Fatalf("server error: %v", err)
	}

	// keep imports used before their activation step
	_ = context.Background
	_ = os.Stderr
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
