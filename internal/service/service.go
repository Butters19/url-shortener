package service

import (
	"errors"

	"github.com/Butters19/url-shortener/internal/generator"
	"github.com/Butters19/url-shortener/internal/storage"
)

const maxRetries = 5

var ErrNotFound = errors.New("url not found")
var ErrInternal = errors.New("internal error")

type Service struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Shorten(originalURL string) (string, error) {
	existingCode, err := s.storage.GetByURL(originalURL)
	if err == nil {
		return existingCode, nil
	}

	for range maxRetries {
		code, err := generator.Generate()
		if err != nil {
			return "", ErrInternal
		}

		err = s.storage.Save(originalURL, code)
		if err == nil {
			return code, nil
		}
		if !errors.Is(err, storage.ErrAlreadyExists) {
			return "", ErrInternal
		}
	}

	return "", ErrInternal
}

func (s *Service) Resolve(shortCode string) (string, error) {
	url, err := s.storage.GetByCode(shortCode)
	if errors.Is(err, storage.ErrNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", ErrInternal
	}
	return url, nil
}