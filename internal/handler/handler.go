package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Butters19/url-shortener/internal/generator"
	"github.com/Butters19/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

const maxRetries = 5

type Handler struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.handleCreate)
	r.Get("/{code}", h.handleGet)
	return r
}

// POST / — принимает {"url": "https://..."}, возвращает {"short_code": "..."}
func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Если URL уже существует — возвращаем существующий код
	existingCode, err := h.storage.GetByURL(req.URL)
	if err == nil {
		writeJSON(w, map[string]string{"short_code": existingCode}, http.StatusOK)
		return
	}

	// Генерируем уникальный код с retry на случай коллизии
	var shortCode string
	for range maxRetries {
		code, err := generator.Generate()
		if err != nil {
			writeError(w, "failed to generate code", http.StatusInternalServerError)
			return
		}

		err = h.storage.Save(req.URL, code)
		if err == nil {
			shortCode = code
			break
		}
		if !errors.Is(err, storage.ErrAlreadyExists) {
			writeError(w, "failed to save url", http.StatusInternalServerError)
			return
		}
	}

	if shortCode == "" {
		writeError(w, "failed to generate unique code", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"short_code": shortCode}, http.StatusCreated)
}

// GET /{code} — возвращает {"url": "https://..."}
func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	originalURL, err := h.storage.GetByCode(code)
	if errors.Is(err, storage.ErrNotFound) {
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
