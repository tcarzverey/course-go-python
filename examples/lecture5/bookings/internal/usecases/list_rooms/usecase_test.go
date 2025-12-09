package list_rooms

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tcarzverey/bookings/internal/models"
	"github.com/tcarzverey/bookings/internal/repository/rooms"
	"go.uber.org/mock/gomock"
)

func TestUsecase_ListRooms(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		roomsRepo func(*MockRoomsRepository)
		want      []models.Room
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			roomsRepo: func(r *MockRoomsRepository) {
				r.EXPECT().ListRooms(ctx).Return([]rooms.Room{
					{ID: 1, Capacity: 10, Name: "Переговорка 1"},
					{ID: 2, Capacity: 15, Name: "Переговорка 2"},
					{ID: 3, Capacity: 3, Name: "Переговорка 3"},
				}, nil)
			},
			want: []models.Room{
				{Id: 1, Capacity: 10, Name: "Переговорка 1"},
				{Id: 2, Capacity: 15, Name: "Переговорка 2"},
				{Id: 3, Capacity: 3, Name: "Переговорка 3"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success: empty",
			roomsRepo: func(r *MockRoomsRepository) {
				r.EXPECT().ListRooms(ctx).Return(nil, nil)
			},
			want:    []models.Room{},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			roomsRepo: func(r *MockRoomsRepository) {
				r.EXPECT().ListRooms(ctx).Return(nil, fmt.Errorf("error"))
			},
			wantErr: assert.Error,
		},
		{
			name: "error: complicated check",
			roomsRepo: func(r *MockRoomsRepository) {
				r.EXPECT().ListRooms(ctx).Return(nil, fmt.Errorf("error"))
			},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.ErrorContains(t, err, "roomsRepo.ListRooms", msgAndArgs...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			roomsMock := NewMockRoomsRepository(ctrl)
			if tt.roomsRepo != nil {
				tt.roomsRepo(roomsMock)
			}

			u := NewUsecase(roomsMock)
			got, err := u.ListRooms(ctx)
			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}
