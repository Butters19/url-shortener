package storage

import "errors"

// ErrNotFound возвращается когда ссылка не найдена
var ErrNotFound = errors.New("url not found")

// ErrAlreadyExists возвращается когда оригинальный URL уже есть в хранилище
var ErrAlreadyExists = errors.New("url already exists")

// Storage — интерфейс для работы с хранилищем ссылок
type Storage interface {
	// Save сохраняет пару оригинальный URL → короткий код
	// Возвращает ErrAlreadyExists если такой URL уже сохранён
	Save(originalURL, shortCode string) error

	// GetByCode возвращает оригинальный URL по короткому коду
	// Возвращает ErrNotFound если код не найден
	GetByCode(shortCode string) (string, error)

	// GetByURL возвращает короткий код по оригинальному URL
	// Возвращает ErrNotFound если URL не найден
	GetByURL(originalURL string) (string, error)
}
