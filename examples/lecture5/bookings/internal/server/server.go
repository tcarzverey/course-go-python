package server

import (
	"context"
	"log/slog"

	"github.com/AlekSi/pointer"
	"github.com/tcarzverey/bookings/internal/generated/api"
)

type Server struct {
	listRoomsUC ListRoomsUsecase
}

func New(listRoomsUC ListRoomsUsecase) *Server {
	return &Server{
		listRoomsUC: listRoomsUC,
	}
}

func (s Server) ListBookings(ctx context.Context, request api.ListBookingsRequestObject) (api.ListBookingsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) CreateBooking(ctx context.Context, request api.CreateBookingRequestObject) (api.CreateBookingResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) CancelBooking(ctx context.Context, request api.CancelBookingRequestObject) (api.CancelBookingResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) ListRooms(ctx context.Context, request api.ListRoomsRequestObject) (api.ListRoomsResponseObject, error) {
	rooms, err := s.listRoomsUC.ListRooms(ctx)
	if err != nil {
		slog.Error("list rooms error")
		return &api.ListRooms500JSONResponse{N500JSONResponse: api.N500JSONResponse{
			LogicCode: pointer.To("listRoomsInternalError"),
			Message:   "error while listing rooms",
			Details: &map[string]any{
				"error": err.Error(),
			},
		}}, nil
	}

	return api.ListRooms200JSONResponse(convertRooms(rooms)), nil
}

func (s Server) CreateRoom(ctx context.Context, request api.CreateRoomRequestObject) (api.CreateRoomResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) ListFreeSlots(ctx context.Context, request api.ListFreeSlotsRequestObject) (api.ListFreeSlotsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

var _ api.StrictServerInterface = (*Server)(nil)
