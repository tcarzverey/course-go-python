package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	// step7 "github.com/observability-practicum/stats-service/internal/telemetry"
	// step7 "go.opentelemetry.io/contrib/bridges/otelslog"
	// step7 "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/observability-practicum/stats-service/internal/handler"
)

func main() {
	port := envOr("PORT", "8081")
	// step7 otlpEndpoint := envOr("OTEL_EXPORTER_OTLP_ENDPOINT", "otel-collector:4317")

	// ── Logger ───────────────────────────────────────────────────────────────
	logger := slog.Default()

	// step7 ctx := context.Background()
	// step7 shutdown, err := telemetry.Init(ctx, "stats-service", otlpEndpoint)
	// step7 if err != nil {
	// step7 	log.Fatalf("failed to initialise OTel SDK: %v", err)
	// step7 }
	// step7 defer func() { _ = shutdown(context.Background()) }()

	// step7 logger = slog.New(otelslog.NewHandler("stats-service"))
	// step7 slog.SetDefault(logger)

	// ── Handler & Router ──────────────────────────────────────────────────────
	h := handler.New(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /track", h.Track)
	mux.HandleFunc("GET /stats/{code}", h.Stats)
	mux.HandleFunc("GET /health", h.Health)

	var root http.Handler = mux
	// step7 root = otelhttp.NewHandler(root, "stats-service-http")

	// ── Start server ──────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", port)
	log.Printf("starting stats-service on %s", addr)
	if err := http.ListenAndServe(addr, root); err != nil {
		log.Fatalf("server error: %v", err)
	}

	_ = context.Background
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
