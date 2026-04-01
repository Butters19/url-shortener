package memory_test

import (
	"context"
	"testing"

	"github.com/Butters19/url-shortener/internal/model"
	"github.com/Butters19/url-shortener/internal/storage"
	"github.com/Butters19/url-shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave_And_GetByCode(t *testing.T) {
	s := memory.New()
	ctx := context.Background()

	err := s.Save(ctx, model.URL{Original: "https://ozon.ru", Code: "abc123defg"})
	require.NoError(t, err)

	u, err := s.GetByCode(ctx, "abc123defg")
	require.NoError(t, err)
	assert.Equal(t, "https://ozon.ru", u.Original)
	assert.Equal(t, "abc123defg", u.Code)
}

func TestSave_And_GetByOrigin(t *testing.T) {
	s := memory.New()
	ctx := context.Background()

	err := s.Save(ctx, model.URL{Original: "https://ozon.ru", Code: "abc123defg"})
	require.NoError(t, err)

	u, err := s.GetByOrigin(ctx, "https://ozon.ru")
	require.NoError(t, err)
	assert.Equal(t, "abc123defg", u.Code)
}

func TestSave_DuplicateURL_ReturnsError(t *testing.T) {
	s := memory.New()
	ctx := context.Background()

	err := s.Save(ctx, model.URL{Original: "https://ozon.ru", Code: "abc123defg"})
	require.NoError(t, err)

	err = s.Save(ctx, model.URL{Original: "https://ozon.ru", Code: "xyz999abcd"})
	assert.ErrorIs(t, err, storage.ErrAlreadyExists)
}

func TestGetByCode_NotFound_ReturnsError(t *testing.T) {
	s := memory.New()
	ctx := context.Background()

	_, err := s.GetByCode(ctx, "notexist")
	assert.ErrorIs(t, err, storage.ErrNotFound)
}

func TestGetByOrigin_NotFound_ReturnsError(t *testing.T) {
	s := memory.New()
	ctx := context.Background()

	_, err := s.GetByOrigin(ctx, "https://notexist.com")
	assert.ErrorIs(t, err, storage.ErrNotFound)
}
