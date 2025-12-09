package list_rooms

import (
	"context"
	"fmt"

	"github.com/tcarzverey/bookings/internal/models"
)

type Usecase struct {
	roomsRepo RoomsRepository
}

func NewUsecase(roomsRepo RoomsRepository) *Usecase {
	return &Usecase{
		roomsRepo: roomsRepo,
	}
}

func (u *Usecase) ListRooms(ctx context.Context) ([]models.Room, error) {
	rooms, err := u.roomsRepo.ListRooms(ctx)
	if err != nil {
		return nil, fmt.Errorf("roomsRepo.ListRooms error: %w", err)
	}

	return convertRooms(rooms), nil
}
