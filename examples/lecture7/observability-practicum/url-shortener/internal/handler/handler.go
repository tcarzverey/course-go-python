package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"

	"github.com/observability-practicum/url-shortener/internal/storage"
)

// Handler holds all dependencies for HTTP handlers.
type Handler struct {
	storage          storage.Store
	baseURL          string
	logger           *slog.Logger
	statsServiceURL  string // step7: URL of the stats-service (empty = disabled)
	shortenedTotal   metric.Int64Counter
	redirectDuration metric.Float64Histogram
}

// New creates a new Handler.
// Pass slog.Default() for logger and "" for statsServiceURL until later steps.
func New(s storage.Store, baseURL string, logger *slog.Logger, statsServiceURL string) *Handler {
	meter := otel.GetMeterProvider().Meter("url-shortener")
	shortenedTotal, _ := meter.Int64Counter("urls.shortened",
		metric.WithDescription("Total number of URLs shortened"),
	)
	redirectDuration, _ := meter.Float64Histogram("urls.redirect.duration",
		metric.WithDescription("Duration of the redirect handler"),
		metric.WithUnit("s"),
	)
	return &Handler{
		storage:          s,
		baseURL:          baseURL,
		logger:           logger,
		statsServiceURL:  statsServiceURL,
		shortenedTotal:   shortenedTotal,
		redirectDuration: redirectDuration,
	}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

// Shorten handles POST /shorten
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("url-shortener").Start(ctx, "handler.Shorten")
	defer span.End()

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid request body")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	code, err := h.storage.Save(ctx, req.URL)
	if err != nil {
		h.shortenedTotal.Add(ctx, 1, metric.WithAttributes(attribute.String("status", "error")))
		h.logger.ErrorContext(r.Context(), "failed to save URL", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "storage error")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	h.shortenedTotal.Add(ctx, 1, metric.WithAttributes(attribute.String("status", "success")))

	span.SetAttributes(
		attribute.String("url.code", code),
		attribute.String("url.original", req.URL),
	)
	h.logger.InfoContext(r.Context(), "URL shortened", "code", code, "original_url", req.URL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ShortenResponse{
		Code:     code,
		ShortURL: fmt.Sprintf("%s/%s", h.baseURL, code),
	})
}

// Redirect handles GET /{code}
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	ctx := r.Context()

	ctx, span := otel.Tracer("url-shortener").Start(ctx, "handler.Redirect")
	defer span.End()

	start := time.Now()
	defer func() { h.redirectDuration.Record(ctx, time.Since(start).Seconds()) }()
	span.SetAttributes(attribute.String("url.code", code))

	record, err := h.storage.Get(ctx, code)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			h.logger.WarnContext(ctx, "code not found", "code", code)
			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_ = h.storage.IncrementClicks(ctx, code)

	// Synchronous call to stats-service so the distributed trace is visible end-to-end.
	if h.statsServiceURL != "" {
		if err := h.trackClick(ctx, code); err != nil {
			h.logger.WarnContext(ctx, "failed to track click", "error", err, "code", code)
		}
	}

	h.logger.InfoContext(ctx, "redirect", "code", code, "url", record.OriginalURL)

	http.Redirect(w, r, record.OriginalURL, http.StatusFound)
}

// Stats handles GET /stats/{code}
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	ctx := r.Context()

	record, err := h.storage.Get(ctx, code)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"code":         record.Code,
		"original_url": record.OriginalURL,
		"clicks":       record.Clicks,
		"created_at":   record.CreatedAt,
	})
}

// Health handles GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// trackClick sends a click event to stats-service, propagating the OTel trace context
// via W3C traceparent/tracestate HTTP headers.
func (h *Handler) trackClick(ctx context.Context, code string) error {
	ctx, span := otel.Tracer("url-shortener").Start(ctx, "handler.trackClick")
	defer span.End()

	body, _ := json.Marshal(map[string]string{"code": code})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.statsServiceURL+"/track", bytes.NewReader(body))
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Inject W3C Trace Context headers — this is what makes it a *distributed* trace.
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// keep imports used at all steps
var (
	_ = bytes.NewReader
	_ = context.Background
)
