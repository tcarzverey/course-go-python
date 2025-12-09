package server

import (
	"context"

	"github.com/tcarzverey/bookings/internal/models"
)

type ListRoomsUsecase interface {
	ListRooms(ctx context.Context) ([]models.Room, error)
}
