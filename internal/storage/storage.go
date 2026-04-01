package storage

import (
	"context"
	"errors"

	"github.com/Butters19/url-shortener/internal/model"
)

var ErrNotFound = errors.New("url not found")
var ErrAlreadyExists = errors.New("url already exists")

type Storage interface {
	Save(ctx context.Context, url model.URL) error
	GetByCode(ctx context.Context, code string) (*model.URL, error)
	GetByOrigin(ctx context.Context, origin string) (*model.URL, error)
}