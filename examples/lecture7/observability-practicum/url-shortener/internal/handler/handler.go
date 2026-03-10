package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	// step7 "go.opentelemetry.io/otel/propagation"

	"github.com/observability-practicum/url-shortener/internal/storage"
)

// Handler holds all dependencies for HTTP handlers.
type Handler struct {
	storage         storage.Store
	baseURL         string
	logger          *slog.Logger
	statsServiceURL string // step7: URL of the stats-service (empty = disabled)
}

// New creates a new Handler.
// Pass slog.Default() for logger and "" for statsServiceURL until later steps.
func New(s storage.Store, baseURL string, logger *slog.Logger, statsServiceURL string) *Handler {
	return &Handler{
		storage:         s,
		baseURL:         baseURL,
		logger:          logger,
		statsServiceURL: statsServiceURL,
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

	code, err := h.storage.Save(req.URL)
	if err != nil {
		// step5 h.logger.ErrorContext(r.Context(), "failed to save URL", "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, "storage error")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	span.SetAttributes(
		attribute.String("url.code", code),
		attribute.String("url.original", req.URL),
	)
	// step5 h.logger.InfoContext(r.Context(), "URL shortened", "code", code, "original_url", req.URL)

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
	_ = ctx // used by step5 logger calls and step3/step7; suppresses "declared and not used" before those steps

	ctx, span := otel.Tracer("url-shortener").Start(ctx, "handler.Redirect")
	defer span.End()
	span.SetAttributes(attribute.String("url.code", code))

	record, err := h.storage.Get(code)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			// step5 h.logger.WarnContext(ctx, "code not found", "code", code)
			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_ = h.storage.IncrementClicks(code)

	// step7 // Synchronous call to stats-service so the distributed trace is visible end-to-end.
	// step7 if h.statsServiceURL != "" {
	// step7 	if err := h.trackClick(ctx, code); err != nil {
	// step7 		h.logger.WarnContext(ctx, "failed to track click", "error", err, "code", code)
	// step7 	}
	// step7 }

	// step5 h.logger.InfoContext(ctx, "redirect", "code", code, "url", record.OriginalURL)

	http.Redirect(w, r, record.OriginalURL, http.StatusFound)
}

// Stats handles GET /stats/{code}
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	record, err := h.storage.Get(code)
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

// step7 // trackClick sends a click event to stats-service, propagating the OTel trace context
// step7 // via W3C traceparent/tracestate HTTP headers.
// step7 func (h *Handler) trackClick(ctx context.Context, code string) error {
// step7 	ctx, span := otel.Tracer("url-shortener").Start(ctx, "handler.trackClick")
// step7 	defer span.End()
// step7
// step7 	body, _ := json.Marshal(map[string]string{"code": code})
// step7 	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.statsServiceURL+"/track", bytes.NewReader(body))
// step7 	if err != nil {
// step7 		span.RecordError(err)
// step7 		return fmt.Errorf("build request: %w", err)
// step7 	}
// step7 	req.Header.Set("Content-Type", "application/json")
// step7
// step7 	// Inject W3C Trace Context headers — this is what makes it a *distributed* trace.
// step7 	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
// step7
// step7 	resp, err := http.DefaultClient.Do(req)
// step7 	if err != nil {
// step7 		span.RecordError(err)
// step7 		return fmt.Errorf("do request: %w", err)
// step7 	}
// step7 	defer resp.Body.Close()
// step7 	return nil
// step7 }

// keep imports used at all steps
var (
	_ = bytes.NewReader
	_ = context.Background
)
