package storage

import (
	"context"
	"errors"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// PGXStorage is a PostgreSQL-backed Store using pgx/v5.
type PGXStorage struct {
	pool *pgxpool.Pool
}

// NewPGX creates a connection pool and ensures the urls table exists.
func NewPGX(ctx context.Context, dsn string) (Store, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// otelpgx instruments every SQL query as a child span automatically.
	config.ConnConfig.Tracer = otelpgx.NewTracer()

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

	meter := otel.GetMeterProvider().Meter("url-shortener")
	_, _ = meter.Int64ObservableGauge("db.pool.connections",
		metric.WithDescription("Number of connections in the pgx connection pool"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			stat := pool.Stat()
			o.Observe(int64(stat.AcquiredConns()), metric.WithAttributes(attribute.String("state", "acquired")))
			o.Observe(int64(stat.IdleConns()), metric.WithAttributes(attribute.String("state", "idle")))
			o.Observe(int64(stat.TotalConns()), metric.WithAttributes(attribute.String("state", "total")))
			return nil
		}),
	)

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

func (s *PGXStorage) Save(ctx context.Context, originalURL string) (string, error) {
	code := generateCode()
	_, err := s.pool.Exec(ctx,
		`INSERT INTO urls (code, original_url) VALUES ($1, $2)`,
		code, originalURL,
	)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *PGXStorage) Get(ctx context.Context, code string) (*URLRecord, error) {
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

func (s *PGXStorage) IncrementClicks(ctx context.Context, code string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE urls SET clicks = clicks + 1 WHERE code = $1`,
		code,
	)
	return err
}
