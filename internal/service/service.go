package service

import (
	"context"
	"errors"
	"time"

	"github.com/Butters19/url-shortener/internal/generator"
	"github.com/Butters19/url-shortener/internal/model"
	"github.com/Butters19/url-shortener/internal/storage"
)

const maxRetries = 5

var ErrNotFound = errors.New("url not found")
var ErrInternal = errors.New("internal error")

type Service struct {
	storage storage.Storage
}

func New(s storage.Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) Shorten(ctx context.Context, originalURL string) (string, error) {
	existing, err := s.storage.GetByOrigin(ctx, originalURL)
	if err == nil {
		return existing.Code, nil
	}

	for range maxRetries {
		code, err := generator.Generate()
		if err != nil {
			return "", ErrInternal
		}

		err = s.storage.Save(ctx, model.URL{
			Original:  originalURL,
			Code:      code,
			CreatedAt: time.Now(),
		})
		if err == nil {
			return code, nil
		}
		if !errors.Is(err, storage.ErrAlreadyExists) {
			return "", ErrInternal
		}
	}

	return "", ErrInternal
}

func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	u, err := s.storage.GetByCode(ctx, code)
	if errors.Is(err, storage.ErrNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", ErrInternal
	}
	return u.Original, nil
}
