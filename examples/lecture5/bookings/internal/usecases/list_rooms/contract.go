//go:generate go tool mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}
package list_rooms

import (
	"context"

	"github.com/tcarzverey/bookings/internal/repository/rooms"
)

type RoomsRepository interface {
	ListRooms(ctx context.Context) ([]rooms.Room, error)
}
