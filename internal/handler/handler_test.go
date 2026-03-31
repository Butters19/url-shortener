package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Butters19/url-shortener/internal/handler"
	"github.com/Butters19/url-shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockService — заглушка сервиса для тестов
type mockService struct {
	shortenFunc func(url string) (string, error)
	resolveFunc func(code string) (string, error)
}

func (m *mockService) Shorten(url string) (string, error) {
	return m.shortenFunc(url)
}

func (m *mockService) Resolve(code string) (string, error) {
	return m.resolveFunc(code)
}

func TestHandleCreate_Success(t *testing.T) {
	svc := &mockService{
		shortenFunc: func(url string) (string, error) {
			return "abc123defg", nil
		},
	}

	h := handler.New(svc)
	body := bytes.NewBufferString(`{"url": "https://ozon.ru"}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "abc123defg", resp["short_code"])
}

func TestHandleCreate_EmptyURL_ReturnsBadRequest(t *testing.T) {
	svc := &mockService{
		shortenFunc: func(url string) (string, error) {
			return "", nil
		},
	}

	h := handler.New(svc)
	body := bytes.NewBufferString(`{"url": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreate_InvalidBody_ReturnsBadRequest(t *testing.T) {
	svc := &mockService{
		shortenFunc: func(url string) (string, error) {
			return "", nil
		},
	}

	h := handler.New(svc)
	body := bytes.NewBufferString(`not json`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleGet_Success(t *testing.T) {
	svc := &mockService{
		resolveFunc: func(code string) (string, error) {
			return "https://ozon.ru", nil
		},
	}

	h := handler.New(svc)
	req := httptest.NewRequest(http.MethodGet, "/abc123defg", nil)
	w := httptest.NewRecorder()

	h.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "https://ozon.ru", resp["url"])
}

func TestHandleGet_NotFound_Returns404(t *testing.T) {
	svc := &mockService{
		resolveFunc: func(code string) (string, error) {
			return "", service.ErrNotFound
		},
	}

	h := handler.New(svc)
	req := httptest.NewRequest(http.MethodGet, "/notexist", nil)
	w := httptest.NewRecorder()

	h.Routes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}