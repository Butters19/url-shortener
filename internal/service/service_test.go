package service_test

import (
	"errors"
	"testing"

	"github.com/Butters19/url-shortener/internal/service"
	"github.com/Butters19/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	urlToCode map[string]string
	codeToURL map[string]string
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		urlToCode: make(map[string]string),
		codeToURL: make(map[string]string),
	}
}

func (m *mockStorage) Save(originalURL, shortCode string) error {
	if _, exists := m.urlToCode[originalURL]; exists {
		return storage.ErrAlreadyExists
	}
	m.urlToCode[originalURL] = shortCode
	m.codeToURL[shortCode] = originalURL
	return nil
}

func (m *mockStorage) GetByCode(shortCode string) (string, error) {
	url, exists := m.codeToURL[shortCode]
	if !exists {
		return "", storage.ErrNotFound
	}
	return url, nil
}

func (m *mockStorage) GetByURL(originalURL string) (string, error) {
	code, exists := m.urlToCode[originalURL]
	if !exists {
		return "", storage.ErrNotFound
	}
	return code, nil
}

func TestShorten_ReturnsShortCode(t *testing.T) {
	svc := service.New(newMockStorage())

	code, err := svc.Shorten("https://ozon.ru")
	require.NoError(t, err)
	assert.Len(t, code, 10)
}

func TestShorten_SameURL_ReturnsSameCode(t *testing.T) {
	svc := service.New(newMockStorage())

	code1, err := svc.Shorten("https://ozon.ru")
	require.NoError(t, err)

	code2, err := svc.Shorten("https://ozon.ru")
	require.NoError(t, err)

	assert.Equal(t, code1, code2)
}

func TestShorten_DifferentURLs_ReturnDifferentCodes(t *testing.T) {
	svc := service.New(newMockStorage())

	code1, err := svc.Shorten("https://ozon.ru")
	require.NoError(t, err)

	code2, err := svc.Shorten("https://github.com")
	require.NoError(t, err)

	assert.NotEqual(t, code1, code2)
}

func TestResolve_ReturnsOriginalURL(t *testing.T) {
	svc := service.New(newMockStorage())

	code, err := svc.Shorten("https://ozon.ru")
	require.NoError(t, err)

	url, err := svc.Resolve(code)
	require.NoError(t, err)
	assert.Equal(t, "https://ozon.ru", url)
}

func TestResolve_NotFound_ReturnsError(t *testing.T) {
	svc := service.New(newMockStorage())

	_, err := svc.Resolve("notexist")
	assert.True(t, errors.Is(err, service.ErrNotFound))
}
