package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	// step8 "github.com/exaring/otelpgx"
)

// PGXStorage is a PostgreSQL-backed Store using pgx/v5.
// Activate with: find "// step8 " → replace with "" (see PLAN.md, Step 8).
type PGXStorage struct {
	pool *pgxpool.Pool
}

// NewPGX creates a connection pool and ensures the urls table exists.
func NewPGX(ctx context.Context, dsn string) (Store, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// step8 // otelpgx instruments every SQL query as a child span automatically.
	// step8 config.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	if err := migrate(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}

	return &PGXStorage{pool: pool}, nil
}

func migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			code         TEXT PRIMARY KEY,
			original_url TEXT        NOT NULL,
			created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
			clicks       BIGINT      NOT NULL DEFAULT 0
		)
	`)
	return err
}

func (s *PGXStorage) Save(originalURL string) (string, error) {
	code := generateCode()
	ctx := context.Background()
	_, err := s.pool.Exec(ctx,
		`INSERT INTO urls (code, original_url) VALUES ($1, $2)`,
		code, originalURL,
	)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *PGXStorage) Get(code string) (*URLRecord, error) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx,
		`SELECT code, original_url, created_at, clicks FROM urls WHERE code = $1`,
		code,
	)
	var r URLRecord
	err := row.Scan(&r.Code, &r.OriginalURL, &r.CreatedAt, &r.Clicks)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &r, nil
}

func (s *PGXStorage) IncrementClicks(code string) error {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx,
		`UPDATE urls SET clicks = clicks + 1 WHERE code = $1`,
		code,
	)
	return err
}

