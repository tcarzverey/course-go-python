package rooms

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	// Запускаем контейнер PostgreSQL
	container, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithInitScripts(filepath.Join("..", "..", "..", "migrations", "001_init_tables.sql")),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		tc.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(10*time.Second),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	endpoint, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, endpoint)
	require.NoError(t, err)

	//// создаём таблицу для теста
	//_, err = pool.Exec(ctx, `
	//    create table rooms
	//	(
	//		id         int generated always as identity primary key,
	//		name       text        not null,
	//		capacity   int         not null,
	//		created_at timestamptz not null default now(),
	//		updated_at timestamptz not null default now()
	//	);
	//`)
	//require.NoError(t, err)

	t.Cleanup(pool.Close)

	return pool
}

func TestQueries_ListRooms_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ctx := context.Background()

	pool := setupTestDB(t)

	q := &Queries{db: pool}

	createdAt := time.Now().UTC().Truncate(time.Second)
	updatedAt := createdAt.Add(5 * time.Minute)

	// готовим данные
	_, err := pool.Exec(ctx, `
        INSERT INTO rooms (name, capacity, created_at, updated_at)
        VALUES 
            ('Room A', 10, $1, $2),
            ('Room B', 20, $1, $2);
    `, createdAt, updatedAt)
	require.NoError(t, err)

	got, err := q.ListRooms(ctx)
	require.NoError(t, err)
	require.Len(t, got, 2)

	require.Equal(t, "Room A", got[0].Name)
	require.Equal(t, int32(10), got[0].Capacity)
	require.Equal(t, "Room B", got[1].Name)
	require.Equal(t, int32(20), got[1].Capacity)
}
