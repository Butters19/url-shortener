package postgres

import (
	"context"
	"errors"

	"github.com/Butters19/url-shortener/internal/model"
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

func (s *Storage) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *Storage) Save(ctx context.Context, url model.URL) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO urls (original_url, short_code) VALUES ($1, $2)`,
		url.Original, url.Code,
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

func (s *Storage) GetByCode(ctx context.Context, code string) (*model.URL, error) {
	var u model.URL
	err := s.pool.QueryRow(ctx,
		`SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = $1`,
		code,
	).Scan(&u.ID, &u.Original, &u.Code, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Storage) GetByOrigin(ctx context.Context, origin string) (*model.URL, error) {
	var u model.URL
	err := s.pool.QueryRow(ctx,
		`SELECT id, original_url, short_code, created_at FROM urls WHERE original_url = $1`,
		origin,
	).Scan(&u.ID, &u.Original, &u.Code, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
