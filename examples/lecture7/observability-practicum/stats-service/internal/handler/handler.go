package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

// clickStore holds per-code click counters, thread-safe via atomic.
type clickStore struct {
	mu     sync.RWMutex
	counts map[string]*atomic.Int64
}

func newClickStore() *clickStore {
	return &clickStore{counts: make(map[string]*atomic.Int64)}
}

func (cs *clickStore) inc(code string) int64 {
	cs.mu.Lock()
	c, ok := cs.counts[code]
	if !ok {
		c = &atomic.Int64{}
		cs.counts[code] = c
	}
	cs.mu.Unlock()
	return c.Add(1)
}

func (cs *clickStore) get(code string) int64 {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	if c, ok := cs.counts[code]; ok {
		return c.Load()
	}
	return 0
}

// Handler is the HTTP handler for the stats-service.
type Handler struct {
	store  *clickStore
	logger *slog.Logger
}

func New(logger *slog.Logger) *Handler {
	return &Handler{store: newClickStore(), logger: logger}
}

// Track handles POST /track
// Body: {"code": "abc123"}
// This endpoint is called by url-shortener on every redirect.
func (h *Handler) Track(w http.ResponseWriter, r *http.Request) {
	// ctx is declared here (always). At step7 we *reassign* it (= not :=) with the
	// extracted trace context so this span joins the trace started in url-shortener.
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
	ctx, span := otel.Tracer("stats-service").Start(ctx, "handler.Track")
	defer span.End()

	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	total := h.store.inc(body.Code)

	span.SetAttributes(
		attribute.String("url.code", body.Code),
		attribute.Int64("clicks.total", total),
	)

	h.logger.InfoContext(ctx, "click tracked", "code", body.Code, "total", total)

	w.WriteHeader(http.StatusNoContent)
}

// Stats handles GET /stats/{code}
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	total := h.store.get(code)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"code":   code,
		"clicks": total,
	})
}

// Health handles GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
