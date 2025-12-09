package rooms

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestQueries_ListRooms(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(ctx)

	q := &Queries{db: mock}

	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Minute)

	// Готовим "результат" запроса
	rows := pgxmock.NewRows([]string{
		"id",
		"name",
		"capacity",
		"created_at",
		"updated_at",
	}).AddRow(
		int32(1),
		"Room A",
		int32(10),
		createdAt,
		updatedAt,
	).AddRow(
		int32(2),
		"Room B",
		int32(20),
		createdAt,
		updatedAt,
	)

	mock.ExpectQuery(listRooms).WillReturnRows(rows)

	// Вызываем метод
	got, err := q.ListRooms(ctx)
	require.NoError(t, err)
	require.Len(t, got, 2)

	require.Equal(t, int32(1), got[0].ID)
	require.Equal(t, "Room A", got[0].Name)
	require.Equal(t, int32(10), got[0].Capacity)

	// Проверяем, что все ожидания pgxmock были выполнены
	require.NoError(t, mock.ExpectationsWereMet())
}
