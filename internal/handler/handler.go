package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Butters19/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type Service interface {
	Shorten(ctx context.Context, originalURL string) (string, error)
	Resolve(ctx context.Context, code string) (string, error)
}

type Handler struct {
	service Service
}

func New(svc Service) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.handleCreate)
	r.Get("/{code}", h.handleGet)
	return r
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	code, err := h.service.Shorten(r.Context(), req.URL)
	if err != nil {
		writeError(w, "failed to shorten url", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"short_code": code}, http.StatusCreated)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	originalURL, err := h.service.Resolve(r.Context(), code)
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, "url not found", http.StatusNotFound)
		return
	}
	if err != nil {
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"url": originalURL}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, map[string]string{"error": msg}, status)
}
