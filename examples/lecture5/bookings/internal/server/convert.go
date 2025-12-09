package server

import (
	"github.com/tcarzverey/bookings/internal/generated/api"
	"github.com/tcarzverey/bookings/internal/models"
)

func convertRooms(rooms []models.Room) []api.Room {
	res := make([]api.Room, 0, len(rooms))
	for _, room := range rooms {
		res = append(res, api.Room{
			Capacity: &room.Capacity,
			Id:       &room.Id,
			Name:     &room.Name,
		})
	}

	return res
}
