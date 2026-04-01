package memory

import (
	"context"
	"sync"

	"github.com/Butters19/url-shortener/internal/model"
	"github.com/Butters19/url-shortener/internal/storage"
)

type Storage struct {
	mu            sync.RWMutex
	originToURL   map[string]*model.URL
	codeToURL     map[string]*model.URL
}

func New() *Storage {
	return &Storage{
		originToURL: make(map[string]*model.URL),
		codeToURL:   make(map[string]*model.URL),
	}
}

func (s *Storage) Save(ctx context.Context, url model.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.originToURL[url.Original]; exists {
		return storage.ErrAlreadyExists
	}

	u := &url
	s.originToURL[url.Original] = u
	s.codeToURL[url.Code] = u
	return nil
}

func (s *Storage) GetByCode(ctx context.Context, code string) (*model.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, exists := s.codeToURL[code]
	if !exists {
		return nil, storage.ErrNotFound
	}
	return u, nil
}

func (s *Storage) GetByOrigin(ctx context.Context, origin string) (*model.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, exists := s.originToURL[origin]
	if !exists {
		return nil, storage.ErrNotFound
	}
	return u, nil
}
