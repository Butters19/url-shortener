package memory_test

import (
	"testing"

	"github.com/Butters19/url-shortener/internal/storage"
	"github.com/Butters19/url-shortener/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave_And_GetByCode(t *testing.T) {
	s := memory.New()

	err := s.Save("https://google.com", "abc123")
	require.NoError(t, err)

	url, err := s.GetByCode("abc123")
	require.NoError(t, err)
	assert.Equal(t, "https://google.com", url)
}

func TestSave_And_GetByURL(t *testing.T) {
	s := memory.New()

	err := s.Save("https://google.com", "abc123")
	require.NoError(t, err)

	code, err := s.GetByURL("https://google.com")
	require.NoError(t, err)
	assert.Equal(t, "abc123", code)
}

func TestSave_DuplicateURL_ReturnsError(t *testing.T) {
	s := memory.New()

	err := s.Save("https://google.com", "abc123")
	require.NoError(t, err)

	err = s.Save("https://google.com", "xyz999")
	assert.ErrorIs(t, err, storage.ErrAlreadyExists)
}

func TestGetByCode_NotFound_ReturnsError(t *testing.T) {
	s := memory.New()

	_, err := s.GetByCode("notexist")
	assert.ErrorIs(t, err, storage.ErrNotFound)
}

func TestGetByURL_NotFound_ReturnsError(t *testing.T) {
	s := memory.New()

	_, err := s.GetByURL("https://notexist.com")
	assert.ErrorIs(t, err, storage.ErrNotFound)
}