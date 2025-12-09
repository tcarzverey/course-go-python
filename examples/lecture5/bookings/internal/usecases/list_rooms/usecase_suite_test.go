package list_rooms

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tcarzverey/bookings/internal/models"
	"github.com/tcarzverey/bookings/internal/repository/rooms"
	"go.uber.org/mock/gomock"
)

type UsecaseListRoomsTestSuite struct {
	suite.Suite
	ctrl      *gomock.Controller
	roomsRepo *MockRoomsRepository
	usecase   *Usecase
	ctx       context.Context
}

func (s *UsecaseListRoomsTestSuite) TestSuccess() {
	s.roomsRepo.EXPECT().ListRooms(s.ctx).Return([]rooms.Room{
		{ID: 1, Capacity: 10, Name: "Переговорка 1"},
		{ID: 2, Capacity: 15, Name: "Переговорка 2"},
		{ID: 3, Capacity: 3, Name: "Переговорка 3"},
	}, nil)

	got, err := s.usecase.ListRooms(s.ctx)
	expected := []models.Room{
		{Id: 1, Capacity: 10, Name: "Переговорка 1"},
		{Id: 2, Capacity: 15, Name: "Переговорка 2"},
		{Id: 3, Capacity: 3, Name: "Переговорка 3"},
	}
	s.Assert().NoError(err)
	s.Assert().Equal(expected, got)
}

func (s *UsecaseListRoomsTestSuite) TestSuccessEmpty() {
	s.roomsRepo.EXPECT().ListRooms(s.ctx).Return(nil, nil)

	got, err := s.usecase.ListRooms(s.ctx)
	s.Assert().NoError(err)
	s.Assert().EqualValues([]models.Room{}, got)
}

func (s *UsecaseListRoomsTestSuite) SetupSuite() {
	s.ctx = context.Background()
}

func (s *UsecaseListRoomsTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.roomsRepo = NewMockRoomsRepository(s.ctrl)
	s.usecase = NewUsecase(s.roomsRepo)
}

func (s *UsecaseListRoomsTestSuite) TearDownTest() {
	s.ctrl.Finish()
	s.ctrl = nil
	s.roomsRepo = nil
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseListRoomsTestSuite))
}
