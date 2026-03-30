package memory

import (
	"sync"

	"github.com/Butters19/url-shortener/internal/storage"
)

type Storage struct {
	mu        sync.RWMutex
	urlToCode map[string]string
	codeToURL map[string]string
}

func New() *Storage {
	return &Storage{
		urlToCode: make(map[string]string),
		codeToURL: make(map[string]string),
	}
}

func (s *Storage) Save(originalURL, shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.urlToCode[originalURL]; exists {
		return storage.ErrAlreadyExists
	}

	s.urlToCode[originalURL] = shortCode
	s.codeToURL[shortCode] = originalURL
	return nil
}

func (s *Storage) GetByCode(shortCode string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.codeToURL[shortCode]
	if !exists {
		return "", storage.ErrNotFound
	}
	return url, nil
}

func (s *Storage) GetByURL(originalURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	code, exists := s.urlToCode[originalURL]
	if !exists {
		return "", storage.ErrNotFound
	}
	return code, nil
}
