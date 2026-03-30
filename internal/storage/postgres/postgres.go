package postgres

import (
	"context"
	"errors"

	"github.com/Butters19/url-shortener/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolation = "23505"

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) Save(originalURL, shortCode string) error {
	_, err := s.pool.Exec(
		context.Background(),
		`INSERT INTO urls (original_url, short_code) VALUES ($1, $2)`,
		originalURL, shortCode,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return storage.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (s *Storage) GetByCode(shortCode string) (string, error) {
	var originalURL string
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT original_url FROM urls WHERE short_code = $1`,
		shortCode,
	).Scan(&originalURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (s *Storage) GetByURL(originalURL string) (string, error) {
	var shortCode string
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT short_code FROM urls WHERE original_url = $1`,
		originalURL,
	).Scan(&shortCode)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return shortCode, nil
}

func (s *Storage) Pool() *pgxpool.Pool {
	return s.pool
}
