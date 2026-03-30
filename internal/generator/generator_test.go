package generator_test

import (
	"testing"

	"github.com/Butters19/url-shortener/internal/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_Length(t *testing.T) {
	code, err := generator.Generate()
	require.NoError(t, err)
	assert.Len(t, code, 10)
}

func TestGenerate_ValidCharacters(t *testing.T) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	validChars := make(map[rune]bool)
	for _, c := range alphabet {
		validChars[c] = true
	}

	for range 100 {
		code, err := generator.Generate()
		require.NoError(t, err)
		for _, c := range code {
			assert.True(t, validChars[c], "unexpected character: %c", c)
		}
	}
}

func TestGenerate_Uniqueness(t *testing.T) {
	codes := make(map[string]bool)
	for range 1000 {
		code, err := generator.Generate()
		require.NoError(t, err)
		assert.False(t, codes[code], "duplicate code generated: %s", code)
		codes[code] = true
	}
}