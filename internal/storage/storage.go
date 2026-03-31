package storage

import "errors"

var ErrNotFound = errors.New("url not found")
var ErrAlreadyExists = errors.New("url already exists")

type Storage interface {
	Save(originalURL, shortCode string) error
	GetByCode(shortCode string) (string, error)
	GetByURL(originalURL string) (string, error)
}