package list_rooms

import (
	"github.com/tcarzverey/bookings/internal/models"
	"github.com/tcarzverey/bookings/internal/repository/rooms"
)

func convertRooms(rooms []rooms.Room) []models.Room {
	res := make([]models.Room, 0, len(rooms))
	for _, room := range rooms {
		res = append(res, models.Room{
			Capacity: int(room.Capacity),
			Id:       int(room.ID),
			Name:     room.Name,
		})
	}

	return res
}
