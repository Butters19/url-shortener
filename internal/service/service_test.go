package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Butters19/url-shortener/internal/model"
	"github.com/Butters19/url-shortener/internal/service"
	"github.com/Butters19/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	originToURL map[string]*model.URL
	codeToURL   map[string]*model.URL
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		originToURL: make(map[string]*model.URL),
		codeToURL:   make(map[string]*model.URL),
	}
}

func (m *mockStorage) Save(ctx context.Context, url model.URL) error {
	if _, exists := m.originToURL[url.Original]; exists {
		return storage.ErrAlreadyExists
	}
	u := &url
	m.originToURL[url.Original] = u
	m.codeToURL[url.Code] = u
	return nil
}

func (m *mockStorage) GetByCode(ctx context.Context, code string) (*model.URL, error) {
	u, exists := m.codeToURL[code]
	if !exists {
		return nil, storage.ErrNotFound
	}
	return u, nil
}

func (m *mockStorage) GetByOrigin(ctx context.Context, origin string) (*model.URL, error) {
	u, exists := m.originToURL[origin]
	if !exists {
		return nil, storage.ErrNotFound
	}
	return u, nil
}

func TestShorten_ReturnsShortCode(t *testing.T) {
	svc := service.New(newMockStorage())

	code, err := svc.Shorten(context.Background(), "https://ozon.ru")
	require.NoError(t, err)
	assert.Len(t, code, 10)
}

func TestShorten_SameURL_ReturnsSameCode(t *testing.T) {
	svc := service.New(newMockStorage())

	code1, err := svc.Shorten(context.Background(), "https://ozon.ru")
	require.NoError(t, err)

	code2, err := svc.Shorten(context.Background(), "https://ozon.ru")
	require.NoError(t, err)

	assert.Equal(t, code1, code2)
}

func TestShorten_DifferentURLs_ReturnDifferentCodes(t *testing.T) {
	svc := service.New(newMockStorage())

	code1, err := svc.Shorten(context.Background(), "https://ozon.ru")
	require.NoError(t, err)

	code2, err := svc.Shorten(context.Background(), "https://github.com")
	require.NoError(t, err)

	assert.NotEqual(t, code1, code2)
}

func TestResolve_ReturnsOriginalURL(t *testing.T) {
	svc := service.New(newMockStorage())

	code, err := svc.Shorten(context.Background(), "https://ozon.ru")
	require.NoError(t, err)

	url, err := svc.Resolve(context.Background(), code)
	require.NoError(t, err)
	assert.Equal(t, "https://ozon.ru", url)
}

func TestResolve_NotFound_ReturnsError(t *testing.T) {
	svc := service.New(newMockStorage())

	_, err := svc.Resolve(context.Background(), "notexist")
	assert.True(t, errors.Is(err, service.ErrNotFound))
}
